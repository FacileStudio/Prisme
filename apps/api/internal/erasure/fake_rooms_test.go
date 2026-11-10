package erasure

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (f *fakeSynapse) room(method, path string, r *http.Request) (int, any) {
	roomID, sub, _ := strings.Cut(strings.TrimPrefix(path, "/v1/rooms/"), "/")
	room, exists := f.rooms[roomID]
	if method == http.MethodDelete && sub == "" {
		delete(f.rooms, roomID)
		return http.StatusOK, map[string]any{"kicked_users": []string{}}
	}
	if !exists {
		return http.StatusNotFound, map[string]string{"errcode": "M_NOT_FOUND"}
	}
	switch sub {
	case "members":
		return http.StatusOK, map[string][]string{"members": room.members}
	case "state":
		return http.StatusOK, map[string]any{"state": room.state()}
	case "make_room_admin":
		var body struct {
			UserID string `json:"user_id"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		room.levels[body.UserID] = 100
		return http.StatusOK, map[string]string{}
	}
	return http.StatusNotFound, map[string]string{"errcode": "M_UNRECOGNIZED"}
}

func (r *fakeRoom) state() []map[string]any {
	create := map[string]any{
		"type": "m.room.create", "state_key": "", "sender": r.creator,
		"content": map[string]string{"room_version": r.version},
	}
	levels := map[string]any{
		"type": "m.room.power_levels", "state_key": "", "sender": r.creator,
		"content": map[string]any{"users": r.levels, "events": map[string]int{"m.room.power_levels": 100}},
	}
	return []map[string]any{create, levels}
}

// addRoom creates a room in the fake with its members and power levels.
func (f *fakeSynapse) addRoom(roomID, creator, version string, levels map[string]int64, members ...string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if levels == nil {
		levels = map[string]int64{}
	}
	f.rooms[roomID] = &fakeRoom{creator: creator, version: version, members: members, levels: levels}
}

func (f *fakeSynapse) addUser(id string, media int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.users[id] = false
	f.media[id] = media
}

// read runs fn with the fake's state locked, so a test reads what the
// server goroutine wrote without a race.
func (f *fakeSynapse) read(fn func(*fakeSynapse)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	fn(f)
}

func (f *fakeSynapse) failErase(fail bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.failErasing = fail
}
