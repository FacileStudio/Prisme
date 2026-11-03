package calls

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"time"

	"github.com/FacileStudio/Prisme/apps/api/internal/echo"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
	"gorm.io/gorm/clause"
)

// Receive handles one call event Echo signed. The call ends when Echo
// finishes the room, or when the last person leaves it.
func (s *Service) Receive(ctx context.Context, body []byte, signature string) error {
	if s.echo == nil {
		return errors.Unavailable(ErrNotConfigured)
	}
	if !s.echo.Verify(body, signature) {
		return errors.Unauthorized("invalid Echo signature")
	}
	event, err := decodeEvent(body)
	if err != nil {
		return errors.Invalid("malformed Echo event")
	}
	switch event.Event {
	case "room_finished":
		s.end(ctx, event.ExternalID)
	case "participant_left":
		if s.empty(ctx, event.RoomID) {
			s.end(ctx, event.ExternalID)
		}
	}
	return nil
}

// running reports whether the room has a call with people in it. A call left
// recorded with nobody connected is ended first, so the next one is announced.
func (s *Service) running(ctx context.Context, roomID, echoRoomID string) bool {
	var current schemas.Call
	if err := s.db.WithContext(ctx).Where("room_id = ?", roomID).Take(&current).Error; err != nil {
		return false
	}
	if !s.empty(ctx, echoRoomID) {
		return true
	}
	s.end(ctx, roomID)
	return false
}

func (s *Service) empty(ctx context.Context, echoRoomID string) bool {
	state, err := s.echo.State(ctx, echoRoomID)
	return err == nil && len(state.Participants) == 0
}

// end deletes the room's running call and edits its message into a finished
// call, as its starter. Deleting first makes a repeated event a no-op.
func (s *Service) end(ctx context.Context, roomID string) {
	var ended []schemas.Call
	err := s.db.WithContext(ctx).Clauses(clause.Returning{}).Where("room_id = ?", roomID).Delete(&ended).Error
	if err != nil || len(ended) == 0 {
		return
	}
	call := ended[0]
	if _, err := s.matrix.SendMessage(ctx, call.StartedBy, roomID, endedEdit(call, time.Now())); err != nil {
		s.log.Warn("calls: marking the call ended", slog.String("room_id", roomID), slog.Any("error", err))
	}
}

func decodeEvent(body []byte) (echo.Event, error) {
	var event echo.Event
	err := json.Unmarshal(body, &event)
	return event, err
}

func endedEdit(call schemas.Call, now time.Time) map[string]any {
	duration := now.Sub(call.StartedAt)
	content := map[string]any{
		"msgtype": "m.text",
		"body":    "Call ended after " + spoken(duration),
		"studio.facile.call": CallContent{
			CallID: call.CallID, StartedBy: call.StartedBy, EndedAt: now.UnixMilli(), DurationMS: duration.Milliseconds(),
		},
	}
	return map[string]any{
		"msgtype":       "m.text",
		"body":          "* " + content["body"].(string),
		"m.new_content": content,
		"m.relates_to":  map[string]string{"rel_type": "m.replace", "event_id": call.EventID},
	}
}

func spoken(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "less than a minute"
	case d < time.Hour:
		return plural(int(d.Minutes()), "minute")
	default:
		return plural(int(d.Hours()), "hour") + " " + plural(int(d.Minutes())%60, "minute")
	}
}

func plural(n int, unit string) string {
	if n == 1 {
		return "1 " + unit
	}
	return strconv.Itoa(n) + " " + unit + "s"
}
