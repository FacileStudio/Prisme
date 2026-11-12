package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/testdb"

	"gorm.io/gorm"
)

const authBootstrap = "docker run -d -e POSTGRES_PASSWORD=t -e POSTGRES_DB=t -p 5432:5432 postgres:16"

var (
	authDBOnce sync.Once
	authGorm   *gorm.DB
	authDBErr  error
	authDBCfg  = testdb.Config{Prefix: "prisme_auth_test", Migrate: schemas.Migrate}
)

func openAuthDatabase(t *testing.T) *sql.DB {
	t.Helper()
	url, configured := testdb.URL()
	if !configured {
		testdb.Announce(authBootstrap)
		t.Skip(testdb.SkipReason(authBootstrap))
	}
	authDBOnce.Do(func() { authGorm, authDBErr = testdb.Open(url, authDBCfg) })
	if authDBErr != nil {
		t.Fatalf("test database: %v", authDBErr)
	}
	if err := testdb.Truncate(authGorm, authDBCfg); err != nil {
		t.Fatalf("reset test database: %v", err)
	}
	db, err := authGorm.DB()
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func newTestService(t *testing.T, db *sql.DB, cfg porte.Config) *Service {
	t.Helper()
	return newLoggedService(t, db, cfg, io.Discard)
}

func newLoggedService(t *testing.T, db *sql.DB, cfg porte.Config, out io.Writer) *Service {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(out, nil))
	service, err := NewService(context.Background(), db, cfg, Guards{ServiceKeys: ParseServiceKeys("key-one, key-two")}, logger)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	return service
}

// echoIdentity answers the identity porte put on the request.
func echoIdentity(w http.ResponseWriter, r *http.Request) {
	identity, _ := porte.From(r.Context())
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user_id": identity.UserID, "email": identity.Email, "verified": identity.EmailVerified,
		"actor": identity.Actor, "session": identity.SessionID, "subject": identity.Subject,
	})
}

type exchange struct {
	status      int
	contentType string
	body        string
}

func call(t *testing.T, s *Service, scope string, headers map[string]string) exchange {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, "/api/spaces", nil)
	for name, value := range headers {
		request.Header.Set(name, value)
	}
	recorder := httptest.NewRecorder()
	s.Require(scope)(http.HandlerFunc(echoIdentity)).ServeHTTP(recorder, request)
	return exchange{recorder.Code, recorder.Header().Get("Content-Type"), recorder.Body.String()}
}

func countIn(t *testing.T, db *sql.DB, query string, args ...any) int {
	t.Helper()
	var n int
	if err := db.QueryRow(query, args...).Scan(&n); err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

// The answers a service-key call gets as of v0.13.0, before the acting
// subject header existed. A call that carries no subject header must keep
// getting exactly these bytes.
func TestServiceKeyCallsWithoutASubjectAnswerAsBefore(t *testing.T) {
	db := openAuthDatabase(t)
	s := newTestService(t, db, porte.Config{})
	cases := []struct {
		name    string
		scope   string
		headers map[string]string
		want    exchange
	}{
		{"creates the account", ScopeRead,
			map[string]string{"Authorization": "Bearer key-two", ActingForHeader: " Alice <Alice@Facile.Studio> "},
			exchange{200, "application/json", `{"actor":"service-key","email":"alice@facile.studio","session":0,"subject":"","user_id":1,"verified":true}` + "\n"}},
		{"finds it again", ScopeAgents,
			map[string]string{"Authorization": "Bearer key-one", ActingForHeader: "alice@facile.studio"},
			exchange{200, "application/json", `{"actor":"service-key","email":"alice@facile.studio","session":0,"subject":"","user_id":1,"verified":true}` + "\n"}},
		{"no email", ScopeRead,
			map[string]string{"Authorization": "Bearer key-one"},
			exchange{400, "application/json", `{"error":{"code":"invalid_argument","message":"X-Prisme-Acting-For must carry the user's email"}}` + "\n"}},
		{"not an email", ScopeRead,
			map[string]string{"Authorization": "Bearer key-one", ActingForHeader: "alice"},
			exchange{400, "application/json", `{"error":{"code":"invalid_argument","message":"X-Prisme-Acting-For must carry the user's email"}}` + "\n"}},
		{"wrong key", ScopeRead,
			map[string]string{"Authorization": "Bearer key-three", ActingForHeader: "alice@facile.studio"},
			exchange{401, "application/json", `{"error":{"code":"unauthenticated","message":"invalid session"}}` + "\n"}},
		{"no key", ScopeRead,
			map[string]string{ActingForHeader: "alice@facile.studio"},
			exchange{401, "application/json", `{"error":{"code":"unauthenticated","message":"missing auth"}}` + "\n"}},
		{"session-only route", "",
			map[string]string{"Authorization": "Bearer key-one", ActingForHeader: "alice@facile.studio"},
			exchange{401, "application/json", `{"error":{"code":"unauthenticated","message":"invalid session"}}` + "\n"}},
	}
	for _, tc := range cases {
		if got := call(t, s, tc.scope, tc.headers); got != tc.want {
			t.Errorf("%s: got %+v, want %+v", tc.name, got, tc.want)
		}
	}
	if n := countIn(t, db, `SELECT count(*) FROM porte_users WHERE email = 'alice@facile.studio' AND email_verified`); n != 1 {
		t.Errorf("verified accounts for alice = %d, want 1", n)
	}
	if n := countIn(t, db, `SELECT count(*) FROM porte_identities`); n != 0 {
		t.Errorf("identities = %d, want 0: a service call links no subject", n)
	}
}
