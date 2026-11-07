package spaces

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/porte/directory"

	"gorm.io/gorm"
)

const testIssuer = "https://sso.test"

type erasureFixture struct {
	gone, other int64
	deletions   *Deletions
	erased      func() []string
}

func count(t *testing.T, db *gorm.DB, query string, args ...any) int64 {
	t.Helper()
	var n int64
	if err := db.Raw(query, args...).Scan(&n).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

func seedPerson(t *testing.T, db *gorm.DB, email, subject string) int64 {
	t.Helper()
	var id int64
	err := db.Raw(`INSERT INTO porte_users (email, name) VALUES (?, ?) RETURNING id`, email, email).Scan(&id).Error
	if err != nil {
		t.Fatalf("create user %s: %v", email, err)
	}
	err = db.Exec(`INSERT INTO porte_identities (user_id, provider, subject) VALUES (?, ?, ?)`, id, testIssuer, subject).Error
	if err != nil {
		t.Fatalf("create identity %s: %v", subject, err)
	}
	local := subject
	rows := []any{
		&schemas.Account{UserID: id, Localpart: local, MatrixUserID: "@" + local + ":hs"},
		&schemas.MatrixSession{PorteSessionID: id, UserID: id, SealedToken: "sealed"},
		&schemas.Agent{OwnerID: id, Slug: local + "-bot", Name: "Bot", MatrixUserID: "@" + local + "-bot:hs"},
	}
	for _, row := range rows {
		if err := db.Create(row).Error; err != nil {
			t.Fatalf("seed %T: %v", row, err)
		}
	}
	return id
}

func fakeSynapse(t *testing.T) (*matrix.Client, func() []string) {
	t.Helper()
	var mu sync.Mutex
	var erased []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/_matrix/client/v3/account/deactivate" {
			mu.Lock()
			erased = append(erased, r.URL.Query().Get("user_id"))
			mu.Unlock()
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)
	client := matrix.New(matrix.Config{BaseURL: server.URL, ServerName: "hs", ASToken: "as"})
	return client, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return append([]string(nil), erased...)
	}
}

func seedErasure(t *testing.T, db *gorm.DB) erasureFixture {
	t.Helper()
	gone := seedPerson(t, db, "gone@facile.studio", "gone")
	other := seedPerson(t, db, "other@facile.studio", "other")
	mustExec(t, db, `INSERT INTO spaces (room_id, name, slug, created_by) VALUES ('!a:hs', 'Studio', 'studio', ?)`, gone)
	mustExec(t, db, `INSERT INTO spaces (room_id, name, slug, created_by) VALUES ('!b:hs', 'Mine', 'mine', ?)`, other)
	mustExec(t, db, `INSERT INTO calls (room_id, call_id, event_id, started_by) VALUES ('!a:hs', 'c', 'e', '@gone:hs')`)
	client, erased := fakeSynapse(t)
	service := &Service{db: db, matrix: client}
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	return erasureFixture{gone: gone, other: other, erased: erased,
		deletions: NewDeletions(service, directory.Client{Issuer: testIssuer}, nil, quiet)}
}

func mustExec(t *testing.T, db *gorm.DB, query string, args ...any) {
	t.Helper()
	if err := db.Exec(query, args...).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
}

// A deleted subject loses its Matrix accounts, sessions, agents and porte user,
// the spaces it created stay with no creator, other people are untouched, and a
// second notice changes nothing.
func TestErasePersonDeletesPersonalDataAndKeepsSharedSpaces(t *testing.T) {
	db := openTestDatabase(t)
	f := seedErasure(t, db)
	for range 2 {
		if err := f.deletions.ErasePerson(context.Background(), testIssuer, "gone"); err != nil {
			t.Fatalf("erase: %v", err)
		}
	}
	for _, table := range []string{"accounts", "matrix_sessions", "porte_users", "porte_identities"} {
		column := "user_id"
		if table == "porte_users" {
			column = "id"
		}
		if n := count(t, db, "SELECT count(*) FROM "+table+" WHERE "+column+" = ?", f.gone); n != 0 {
			t.Errorf("%s still holds %d rows of the erased person", table, n)
		}
	}
	if n := count(t, db, `SELECT count(*) FROM agents WHERE owner_id = ?`, f.gone); n != 0 {
		t.Errorf("agents of the erased person remain: %d", n)
	}
	if n := count(t, db, `SELECT count(*) FROM spaces WHERE name = 'Studio' AND created_by = 0`); n != 1 {
		t.Errorf("the shared space lost its row or kept its creator")
	}
	if n := count(t, db, `SELECT count(*) FROM calls WHERE started_by = ''`); n != 1 {
		t.Errorf("the running call kept the erased starter")
	}
	if got := f.erased(); len(got) != 2 {
		t.Errorf("synapse erased %v, want the account and its agent once", got)
	}
}

func TestErasePersonLeavesOthersUntouched(t *testing.T) {
	db := openTestDatabase(t)
	f := seedErasure(t, db)
	if err := f.deletions.ErasePerson(context.Background(), testIssuer, "gone"); err != nil {
		t.Fatalf("erase: %v", err)
	}
	for table, column := range map[string]string{
		"porte_users": "id", "porte_identities": "user_id", "accounts": "user_id",
		"matrix_sessions": "user_id", "agents": "owner_id", "spaces": "created_by",
	} {
		if n := count(t, db, "SELECT count(*) FROM "+table+" WHERE "+column+" = ?", f.other); n != 1 {
			t.Errorf("%s holds %d rows of the other person, want 1", table, n)
		}
	}
	if err := f.deletions.ErasePerson(context.Background(), testIssuer, "nobody"); err != nil {
		t.Fatalf("unknown subject: %v", err)
	}
}

func TestDeletionStateSurvivesARestart(t *testing.T) {
	db := openTestDatabase(t)
	state := deletionState{db: db}
	if saved, err := state.Load(context.Background()); err != nil || saved != nil {
		t.Fatalf("empty Load = %q, %v; want nil", saved, err)
	}
	for _, blob := range []string{"one", "two"} {
		if err := state.Save(context.Background(), []byte(blob)); err != nil {
			t.Fatalf("save: %v", err)
		}
	}
	if saved, err := state.Load(context.Background()); err != nil || string(saved) != "two" {
		t.Fatalf("Load = %q, %v; want two", saved, err)
	}
}
