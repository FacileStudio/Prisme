package push

import (
	"context"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"
)

func TestRecipientsAreTheOtherPersonOrTheMentioned(t *testing.T) {
	direct := event{Sender: "@ada:x"}
	if got := recipients(direct, []string{"@ada:x", "@bob:x"}); !slices.Equal(got, []string{"@bob:x"}) {
		t.Fatalf("direct message recipients = %v, want [@bob:x]", got)
	}
	channel := []string{"@ada:x", "@bob:x", "@cy:x", "@prisme:x"}
	if got := recipients(event{Sender: "@ada:x"}, channel); len(got) != 0 {
		t.Fatalf("channel message without mentions reached %v", got)
	}
	mentioned := event{Sender: "@ada:x", Content: content{Mentions: &mentions{UserIDs: []string{"@cy:x", "@ada:x", "@gone:x"}}}}
	if got := recipients(mentioned, channel); !slices.Equal(got, []string{"@cy:x"}) {
		t.Fatalf("mention recipients = %v, want [@cy:x]", got)
	}
}

func TestOnlyFreshNewMessagesAreNotifiable(t *testing.T) {
	now := time.Now()
	fresh := event{Type: "m.room.message", Timestamp: now.UnixMilli()}
	edit := fresh
	edit.Content.RelatesTo = &relation{RelType: "m.replace"}
	old := fresh
	old.Timestamp = now.Add(-time.Hour).UnixMilli()
	key := ""
	state := fresh
	state.StateKey = &key
	reaction := fresh
	reaction.Type = "m.reaction"
	if !notifiable(fresh, now) {
		t.Fatal("a fresh message is not notifiable")
	}
	for name, e := range map[string]event{"edit": edit, "old": old, "state": state, "reaction": reaction} {
		if notifiable(e, now) {
			t.Errorf("%s is notifiable", name)
		}
	}
}

func TestTransactionsNeedTheHomeserverToken(t *testing.T) {
	s := NewService(Deps{HSToken: "hs"})
	if !s.Authorized("Bearer hs") || s.Authorized("Bearer as") || s.Authorized("hs") {
		t.Fatal("hs_token check is wrong")
	}
	if NewService(Deps{}).Authorized("Bearer ") {
		t.Fatal("an unset hs_token let a request in")
	}
}

func TestFCMReportsAnUnregisteredDeviceAsGone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"status":"NOT_FOUND","details":[{"errorCode":"UNREGISTERED"}]}}`))
	}))
	defer server.Close()
	f := &FCM{endpoint: server.URL, client: server.Client()}
	if err := f.send(context.Background(), "token", message{Title: "t"}); err != errGone {
		t.Fatalf("send = %v, want errGone", err)
	}
}
