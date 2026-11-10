package erasure

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
)

const testServer = "hs"

// fakeRoom is one room the fake homeserver holds.
type fakeRoom struct {
	creator string
	version string
	members []string
	levels  map[string]int64
}

// fakeSynapse is the slice of Synapse's admin and client API the erasure
// uses, with real state, so a test can read what each call changed.
type fakeSynapse struct {
	mu          sync.Mutex
	users       map[string]bool
	rooms       map[string]*fakeRoom
	media       map[string]int
	calls       []string
	failErasing bool
}

func newFakeSynapse(t *testing.T) (*fakeSynapse, *matrix.Client) {
	t.Helper()
	f := &fakeSynapse{users: map[string]bool{}, rooms: map[string]*fakeRoom{}, media: map[string]int{}}
	server := httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(server.Close)
	return f, matrix.New(matrix.Config{BaseURL: server.URL, ServerName: testServer, ASToken: "as"})
}

func (f *fakeSynapse) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	path := strings.TrimPrefix(r.URL.Path, "/_synapse/admin")
	f.calls = append(f.calls, r.Method+" "+path)
	status, body := f.route(r.Method, path, r)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func (f *fakeSynapse) route(method, path string, r *http.Request) (int, any) {
	switch {
	case method == http.MethodPost && path == "/_matrix/client/v3/register":
		return f.register(r)
	case strings.HasPrefix(path, "/v1/users/") || strings.HasPrefix(path, "/v2/users/"):
		return f.user(method, path)
	case strings.HasPrefix(path, "/v1/rooms/"):
		return f.room(method, path, r)
	case method == http.MethodPost && strings.HasPrefix(path, "/v1/deactivate/"):
		return f.deactivate(strings.TrimPrefix(path, "/v1/deactivate/"))
	case path == "/v1/background_updates/enabled":
		return http.StatusOK, map[string]bool{"enabled": true}
	}
	return http.StatusNotFound, map[string]string{"errcode": "M_UNRECOGNIZED"}
}

func (f *fakeSynapse) register(r *http.Request) (int, any) {
	var body struct{ Username string }
	_ = json.NewDecoder(r.Body).Decode(&body)
	id := "@" + body.Username + ":" + testServer
	if _, exists := f.users[id]; exists {
		return http.StatusBadRequest, map[string]string{"errcode": "M_USER_IN_USE"}
	}
	f.users[id] = false
	return http.StatusOK, map[string]string{"user_id": id}
}

func (f *fakeSynapse) user(method, path string) (int, any) {
	rest := path[len("/v1/users/"):]
	id, sub, _ := strings.Cut(rest, "/")
	deactivated, exists := f.users[id]
	switch {
	case sub == "joined_rooms":
		return http.StatusOK, map[string][]string{"joined_rooms": f.roomsOf(id)}
	case !exists:
		return http.StatusNotFound, map[string]string{"errcode": "M_NOT_FOUND"}
	case sub == "media" && method == http.MethodDelete:
		deleted := make([]string, f.media[id])
		f.media[id] = 0
		return http.StatusOK, map[string][]string{"deleted_media": deleted}
	}
	return http.StatusOK, map[string]bool{"deactivated": deactivated}
}

func (f *fakeSynapse) deactivate(id string) (int, any) {
	if f.failErasing {
		return http.StatusServiceUnavailable, map[string]string{"errcode": "M_UNKNOWN"}
	}
	if _, exists := f.users[id]; !exists {
		return http.StatusNotFound, map[string]string{"errcode": "M_NOT_FOUND"}
	}
	f.users[id] = true
	for _, room := range f.rooms {
		room.members = slices.DeleteFunc(room.members, func(m string) bool { return m == id })
	}
	return http.StatusOK, map[string]string{}
}

func (f *fakeSynapse) roomsOf(id string) []string {
	var out []string
	for roomID, room := range f.rooms {
		if slices.Contains(room.members, id) {
			out = append(out, roomID)
		}
	}
	slices.Sort(out)
	return out
}
