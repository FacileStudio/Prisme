package push

import (
	"context"
	stderrors "errors"
	"log/slog"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	maxBodyRunes   = 240
	maxEventAge    = 5 * time.Minute
	deliverTimeout = 30 * time.Second
)

// Receive hands a transaction's new messages to the notifier in the
// background, so Synapse gets its answer at once and never retries.
func (s *Service) Receive(txn transaction) {
	if s.sender == nil {
		return
	}
	now := time.Now()
	events := slices.DeleteFunc(txn.Events, func(e event) bool { return !notifiable(e, now) })
	if len(events) == 0 {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), deliverTimeout)
		defer cancel()
		for _, e := range events {
			s.deliver(ctx, e)
		}
	}()
}

func notifiable(e event, now time.Time) bool {
	if e.Type != "m.room.message" || e.StateKey != nil {
		return false
	}
	if e.Content.RelatesTo != nil && e.Content.RelatesTo.RelType == "m.replace" {
		return false
	}
	return now.Sub(time.UnixMilli(e.Timestamp)) < maxEventAge
}

// recipients is who a message is for: the other person of a two-member room,
// otherwise the members it mentions. The sender is never one of them.
func recipients(e event, members []string) []string {
	others := slices.DeleteFunc(slices.Clone(members), func(id string) bool { return id == e.Sender })
	if len(members) == 2 {
		return others
	}
	if e.Content.Mentions == nil {
		return nil
	}
	return slices.DeleteFunc(others, func(id string) bool { return !slices.Contains(e.Content.Mentions.UserIDs, id) })
}

func (s *Service) deliver(ctx context.Context, e event) {
	members, err := s.people.JoinedMembers(ctx, e.Sender, e.RoomID)
	if err != nil {
		s.deps.Logger.Warn("push: failed to read the room members", slog.Any("error", err))
		return
	}
	targets := recipients(e, members)
	if len(targets) == 0 {
		return
	}
	devices := s.devices(ctx, targets)
	if len(devices) == 0 {
		return
	}
	m := s.compose(ctx, e, len(members) == 2)
	muted := map[string]bool{}
	for _, d := range devices {
		if _, seen := muted[d.MatrixUserID]; !seen {
			muted[d.MatrixUserID] = s.people.RoomMuted(ctx, d.MatrixUserID, e.RoomID)
		}
		if !muted[d.MatrixUserID] {
			s.sendOne(ctx, d.Token, m)
		}
	}
}

func (s *Service) sendOne(ctx context.Context, token string, m message) {
	err := s.sender.send(ctx, token, m)
	if stderrors.Is(err, errGone) {
		s.forget(ctx, token)
	} else if err != nil {
		s.deps.Logger.Warn("push to a device failed", slog.Any("error", err))
	}
}

func (s *Service) compose(ctx context.Context, e event, direct bool) message {
	title := s.people.DisplayName(ctx, e.Sender)
	if !direct {
		if name := s.people.RoomName(ctx, e.Sender, e.RoomID); name != "" {
			title += " · " + name
		}
	}
	return message{Title: title, Body: truncate(e.Content.Body), RoomID: e.RoomID}
}

func truncate(text string) string {
	text = strings.TrimSpace(text)
	if utf8.RuneCountInString(text) <= maxBodyRunes {
		return text
	}
	return string([]rune(text)[:maxBodyRunes-1]) + "…"
}
