package auth

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FacileStudio/porte"
)

const testIssuer = "https://sso.test"

// subjectService is a service whose sign-in issuer is testIssuer, without
// the OIDC discovery a real issuer needs; logs land in the returned buffer.
func subjectService(t *testing.T, db *sql.DB) (*Service, *bytes.Buffer) {
	t.Helper()
	logs := &bytes.Buffer{}
	s := newLoggedService(t, db, porte.Config{}, logs)
	s.services.issuer = testIssuer
	return s, logs
}

func actingAs(email, subject string) map[string]string {
	return map[string]string{"Authorization": "Bearer key-two", ActingForHeader: email, ActingSubjectHeader: subject}
}

func seedAccount(t *testing.T, db *sql.DB, email, provider, subject string) int64 {
	t.Helper()
	var id int64
	if err := db.QueryRow(`INSERT INTO porte_users (email, email_verified) VALUES ($1, true) RETURNING id`, email).Scan(&id); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if subject != "" {
		if _, err := db.Exec(`INSERT INTO porte_identities (user_id, provider, subject, access_token) VALUES ($1, $2, $3, 'kept')`,
			id, provider, subject); err != nil {
			t.Fatalf("seed identity: %v", err)
		}
	}
	return id
}

func identitiesOf(t *testing.T, db *sql.DB) string {
	t.Helper()
	rows, err := db.Query(`SELECT user_id, provider, subject, access_token FROM porte_identities ORDER BY user_id, provider`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	var out []string
	for rows.Next() {
		var id, provider, subject, token string
		if err := rows.Scan(&id, &provider, &subject, &token); err != nil {
			t.Fatal(err)
		}
		out = append(out, id+" "+provider+" "+subject+" "+token)
	}
	return strings.Join(out, "; ")
}

func identityBody(id, email string) string {
	return `{"actor":"service-key","email":"` + email + `","session":0,"subject":"","user_id":` + id + `,"verified":true}` + "\n"
}

// A new person, and an account a service key created by email alone, gain
// the subject as an identity of the issuer, spelled as a sign-in spells it.
// The log names the account and the key, never the email or the subject.
func TestTheActingSubjectLinksAnUnlinkedAccount(t *testing.T) {
	db := openAuthDatabase(t)
	s, logs := subjectService(t, db)
	if got := call(t, s, ScopeRead, map[string]string{"Authorization": "Bearer key-one", ActingForHeader: "old@facile.studio"}); got.status != 200 {
		t.Fatalf("provision by email: %+v", got)
	}
	if got := call(t, s, ScopeRead, actingAs("old@facile.studio", "sub-old")); got.body != identityBody("1", "old@facile.studio") {
		t.Fatalf("link the existing account: %+v", got)
	}
	got := call(t, s, ScopeAgents, actingAs("new@facile.studio", "sub-new"))
	var created string
	if err := db.QueryRow(`SELECT id::text FROM porte_users WHERE email = 'new@facile.studio'`).Scan(&created); err != nil || got.body != identityBody(created, "new@facile.studio") {
		t.Fatalf("create and link: %+v (%v)", got, err)
	}
	if got, want := identitiesOf(t, db), "1 https://sso.test sub-old ; "+created+" https://sso.test sub-new "; got != want {
		t.Errorf("identities = %q, want %q", got, want)
	}
	text := logs.String()
	if strings.Count(text, "linked a registre subject") != 2 || !strings.Contains(text, "user_id="+created) || !strings.Contains(text, "service_key=2") {
		t.Errorf("log does not record each link with the account and key:\n%s", text)
	}
	for _, secret := range []string{"facile.studio", "sub-old", "sub-new"} {
		if strings.Contains(text, secret) {
			t.Errorf("log carries %q:\n%s", secret, text)
		}
	}
}

// An account holding the same subject, under either spelling of the issuer,
// is left exactly as it was: its stored tokens survive.
func TestTheSameSubjectIsLeftAlone(t *testing.T) {
	db := openAuthDatabase(t)
	s, _ := subjectService(t, db)
	seedAccount(t, db, "exact@facile.studio", testIssuer, "sub-exact")
	seedAccount(t, db, "slash@facile.studio", testIssuer+"/", "sub-slash")
	before := identitiesOf(t, db)
	if got := call(t, s, ScopeRead, actingAs("exact@facile.studio", "sub-exact")); got.body != identityBody("1", "exact@facile.studio") {
		t.Fatalf("exact spelling: %+v", got)
	}
	if got := call(t, s, ScopeRead, actingAs("slash@facile.studio", "sub-slash")); got.body != identityBody("2", "slash@facile.studio") {
		t.Fatalf("trailing slash: %+v", got)
	}
	if after := identitiesOf(t, db); after != before {
		t.Errorf("identities changed: %q -> %q", before, after)
	}
}

// Each refusal answers tronc's error envelope with its own stable code and
// writes nothing: no account, no identity.
func TestTheActingSubjectRefusesAndWritesNothing(t *testing.T) {
	db := openAuthDatabase(t)
	s, logs := subjectService(t, db)
	seedAccount(t, db, "held@facile.studio", testIssuer+"/", "sub-held")
	malformed := exchange{400, "application/json", `{"error":{"code":"acting_subject_invalid","message":"X-Prisme-Acting-Subject must carry one Registre subject of 1 to 64 bytes with no control characters"}}` + "\n"}
	taken := exchange{409, "application/json", `{"error":{"code":"acting_subject_taken","message":"this Registre subject belongs to another Prisme account"}}` + "\n"}
	mismatch := exchange{409, "application/json", `{"error":{"code":"acting_subject_mismatch","message":"the Prisme account for this email holds another Registre subject"}}` + "\n"}
	cases := []struct {
		name    string
		email   string
		subject string
		want    exchange
	}{
		{"empty", "fresh@facile.studio", "", malformed},
		{"too long", "fresh@facile.studio", strings.Repeat("a", 65), malformed},
		{"control character", "fresh@facile.studio", "sub\x01x", malformed},
		{"tab", "fresh@facile.studio", "sub\tx", malformed},
		{"not utf-8", "fresh@facile.studio", "sub\xffx", malformed},
		{"subject of another account, new email", "moved@facile.studio", "sub-held", taken},
		{"another subject on the account", "held@facile.studio", "sub-other", mismatch},
	}
	before := identitiesOf(t, db)
	for _, tc := range cases {
		if got := call(t, s, ScopeRead, actingAs(tc.email, tc.subject)); got != tc.want {
			t.Errorf("%s: got %+v, want %+v", tc.name, got, tc.want)
		}
	}
	request := httptest.NewRequest(http.MethodGet, "/api/spaces", nil)
	request.Header.Set("Authorization", "Bearer key-one")
	request.Header.Set(ActingForHeader, "fresh@facile.studio")
	request.Header.Add(ActingSubjectHeader, "sub-a")
	request.Header.Add(ActingSubjectHeader, "sub-b")
	recorder := httptest.NewRecorder()
	s.Require(ScopeRead)(http.HandlerFunc(echoIdentity)).ServeHTTP(recorder, request)
	if recorder.Code != 400 || recorder.Body.String() != malformed.body {
		t.Errorf("two subjects: %d %s", recorder.Code, recorder.Body.String())
	}
	if n := countIn(t, db, `SELECT count(*) FROM porte_users WHERE email <> 'held@facile.studio'`); n != 0 {
		t.Errorf("a refused call created %d accounts", n)
	}
	if after := identitiesOf(t, db); after != before {
		t.Errorf("identities changed: %q -> %q", before, after)
	}
	if n := strings.Count(logs.String(), "refused, nothing written"); n != 2 {
		t.Errorf("logged %d conflicts, want 2:\n%s", n, logs.String())
	}
}

// Subjects of exactly 64 bytes are Registre's longest and are accepted.
func TestASubjectOfRegistresMaximumLengthLinks(t *testing.T) {
	db := openAuthDatabase(t)
	s, _ := subjectService(t, db)
	subject := strings.Repeat("a", 64)
	if got := call(t, s, ScopeRead, actingAs("long@facile.studio", subject)); got.status != 200 {
		t.Fatalf("64-byte subject: %+v", got)
	}
	if n := countIn(t, db, `SELECT count(*) FROM porte_identities WHERE subject = $1`, subject); n != 1 {
		t.Errorf("identities = %d, want 1", n)
	}
}

// Without Registre sign-in there is no issuer to file a subject under, so
// the header is refused rather than silently dropped.
func TestTheActingSubjectNeedsAnIssuer(t *testing.T) {
	db := openAuthDatabase(t)
	s := newTestService(t, db, porte.Config{})
	want := exchange{412, "application/json", `{"error":{"code":"acting_subject_unsupported","message":"X-Prisme-Acting-Subject needs Registre sign-in, which this Prisme does not have"}}` + "\n"}
	if got := call(t, s, ScopeRead, actingAs("fresh@facile.studio", "sub-fresh")); got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if n := countIn(t, db, `SELECT count(*) FROM porte_users`); n != 0 {
		t.Errorf("accounts = %d, want 0", n)
	}
}

// The header is honoured only on service-key calls. A signed-in person who
// sends it, with an acting email too, is themselves and nothing is written.
func TestASessionIgnoresTheActingSubject(t *testing.T) {
	db := openAuthDatabase(t)
	s, logs := subjectService(t, db)
	userID := seedAccount(t, db, "me@facile.studio", "", "")
	if _, err := db.Exec(`INSERT INTO porte_sessions (token_hash, user_id) VALUES ($1, $2)`, porte.HashToken("session-token"), userID); err != nil {
		t.Fatalf("seed session: %v", err)
	}
	headers := map[string]string{"Cookie": porte.SessionCookieName + "=session-token",
		ActingForHeader: "victim@facile.studio", ActingSubjectHeader: "sub-victim"}
	got := call(t, s, ScopeRead, headers)
	if want := `{"actor":"","email":"","session":1,"subject":"","user_id":1,"verified":false}` + "\n"; got.status != 200 || got.body != want {
		t.Fatalf("session call: %+v", got)
	}
	headers[ActingSubjectHeader] = "\x01"
	if again := call(t, s, ScopeRead, headers); again.status != 200 {
		t.Fatalf("a session call failed on a malformed subject header: %+v", again)
	}
	if n := countIn(t, db, `SELECT count(*) FROM porte_identities`); n != 0 {
		t.Errorf("identities = %d, want 0", n)
	}
	if n := countIn(t, db, `SELECT count(*) FROM porte_users`); n != 1 {
		t.Errorf("accounts = %d, want 1", n)
	}
	if logs.Len() != 0 {
		t.Errorf("a session call logged:\n%s", logs.String())
	}
}
