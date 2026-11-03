package calls

import (
	"context"
	"testing"
	"time"

	"github.com/FacileStudio/Prisme/apps/api/internal/echo"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
)

func TestEndedEditReplacesTheCallMessage(t *testing.T) {
	started := time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC)
	call := schemas.Call{RoomID: "!r:prisme.test", CallID: "prisme-1", EventID: "$call", StartedBy: "@alice:prisme.test", StartedAt: started}
	edit := endedEdit(call, started.Add(75*time.Minute))
	relation := edit["m.relates_to"].(map[string]string)
	if relation["rel_type"] != "m.replace" || relation["event_id"] != "$call" {
		t.Fatalf("relation = %v, want an m.replace of $call", relation)
	}
	content := edit["m.new_content"].(map[string]any)
	if content["body"] != "Call ended after 1 hour 15 minutes" {
		t.Errorf("body = %q", content["body"])
	}
	block := content["studio.facile.call"].(CallContent)
	if block.CallID != "prisme-1" || block.EndedAt == 0 || block.DurationMS != (75*time.Minute).Milliseconds() {
		t.Errorf("call block = %+v", block)
	}
}

func TestReceiveRefusesAnUnsignedEvent(t *testing.T) {
	service := &Service{echo: echo.New("http://echo.test", "integration-secret")}
	err := service.Receive(context.Background(), []byte(`{"event":"room_finished"}`), "deadbeef")
	if err == nil {
		t.Fatal("an event with a wrong signature was accepted")
	}
}
