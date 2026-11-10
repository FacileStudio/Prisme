package erasure

import (
	"context"
	"encoding/json"
	"os"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
)

// liveRooms are the rooms a live test creates for the person to erase.
type liveRooms struct {
	solo, direct, shared string
}

// liveSynapse is the homeserver of docker-compose.dev.yml, or any Synapse
// running deploy/synapse's image (whose entrypoint makes the appservice
// sender an admin):
//
//	docker compose -p prisme-dev -f docker-compose.dev.yml up -d --build
//	TEST_SYNAPSE_URL=http://127.0.0.1:18008 TEST_SYNAPSE_AS_TOKEN=dev-as-token \
//	  TEST_SYNAPSE_SERVER_NAME=localhost go test ./internal/erasure -run Live
//
// Without the three variables the live tests skip. Every run registers new
// ids, because Synapse never frees one.
func liveSynapse(t *testing.T) *matrix.Client {
	t.Helper()
	url, token, server := os.Getenv("TEST_SYNAPSE_URL"), os.Getenv("TEST_SYNAPSE_AS_TOKEN"), os.Getenv("TEST_SYNAPSE_SERVER_NAME")
	if url == "" || token == "" || server == "" {
		t.Skip("TEST_SYNAPSE_URL, TEST_SYNAPSE_AS_TOKEN and TEST_SYNAPSE_SERVER_NAME are not set")
	}
	client := matrix.New(matrix.Config{BaseURL: url, ServerName: server, ASToken: token})
	if err := client.AdminReady(context.Background()); err != nil {
		t.Fatalf("the appservice token is no admin credential on this Synapse: %v", err)
	}
	return client
}

func mustLive(t *testing.T, err error, what string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", what, err)
	}
}

// liveWorld registers the person, their agent and a colleague, and builds the
// three rooms the rule tells apart: one only the person and their agent are
// in, a direct chat the person created and alone administers, and a room the
// colleague created.
func liveWorld(t *testing.T, client *matrix.Client, run string) (person, string, liveRooms) {
	ctx := context.Background()
	gone := person{subject: "live-" + run, email: "gone" + run + "@facile.studio"}
	gone.mxid, gone.ag = client.UserID("gone"+run), client.UserID("gone"+run+"-bot")
	colleague := client.UserID("keep" + run)
	for _, local := range []string{"gone" + run, "gone" + run + "-bot", "keep" + run} {
		mustLive(t, client.Register(ctx, local), "register "+local)
	}
	mustLive(t, client.SetProfile(ctx, gone.mxid, "Gone Person", ""), "profile")
	var rooms liveRooms
	var err error
	rooms.solo, err = client.CreateRoom(ctx, gone.mxid, matrix.CreateRoomRequest{Preset: "private_chat", Invite: []string{gone.ag}})
	mustLive(t, err, "create the solo room")
	mustLive(t, client.Join(ctx, gone.ag, rooms.solo), "agent joins")
	rooms.direct, err = client.CreateRoom(ctx, gone.mxid, matrix.CreateRoomRequest{Preset: "private_chat", Invite: []string{colleague}})
	mustLive(t, err, "create the direct chat")
	mustLive(t, client.Join(ctx, colleague, rooms.direct), "colleague joins")
	rooms.shared, err = client.CreateRoom(ctx, colleague, matrix.CreateRoomRequest{Preset: "private_chat", Invite: []string{gone.mxid}})
	mustLive(t, err, "create the shared room")
	mustLive(t, client.Join(ctx, gone.mxid, rooms.shared), "person joins")
	for _, roomID := range []string{rooms.solo, rooms.direct, rooms.shared} {
		_, err := client.SendMessage(ctx, gone.mxid, roomID, map[string]string{"msgtype": "m.text", "body": "personal " + run})
		mustLive(t, err, "send")
	}
	_, err = client.Upload(ctx, gone.mxid, []byte("\x89PNG\r\n\x1a\n"), "image/png")
	mustLive(t, err, "upload")
	return gone, colleague, rooms
}

func liveRun() string {
	return strconv.FormatInt(time.Now().UnixNano()%1_000_000_000, 36)
}

// powerOf reads a user's power level in a room through the admin API.
func powerOf(t *testing.T, client *matrix.Client, roomID, userID string) int64 {
	t.Helper()
	state, err := client.RoomState(context.Background(), roomID)
	mustLive(t, err, "room state")
	for _, event := range state {
		if event.Type == "m.room.power_levels" {
			var content levelsContent
			_ = json.Unmarshal(event.Content, &content)
			return content.Users[userID]
		}
	}
	return 0
}

// waitParted waits for Synapse's background parter to take the user out of
// the room, which it does after the deactivation call returns.
func waitParted(t *testing.T, client *matrix.Client, roomID, userID string) []string {
	t.Helper()
	for range 50 {
		members, err := client.RoomMembers(context.Background(), roomID)
		mustLive(t, err, "members")
		if !slices.Contains(members, userID) {
			return members
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("%s is still in %s", userID, roomID)
	return nil
}
