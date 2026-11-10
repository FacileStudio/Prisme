package erasure

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"testing"

	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/testdb"

	"gorm.io/gorm"
)

const (
	testIssuer    = "https://sso.test"
	testBootstrap = "docker run -d -e POSTGRES_PASSWORD=t -e POSTGRES_DB=t -p 5432:5432 postgres:16"
	testBot       = "@prisme:" + testServer
)

var (
	testDBOnce sync.Once
	testDB     *gorm.DB
	testDBErr  error
	testDBCfg  = testdb.Config{Prefix: "prisme_erasure_test", Migrate: schemas.Migrate}
)

// person is one seeded account and the needles a leftover scan looks for.
type person struct {
	id                       int64
	email, subject, mxid, ag string
}

func openTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	url, configured := testdb.URL()
	if !configured {
		testdb.Announce(testBootstrap)
		t.Skip(testdb.SkipReason(testBootstrap))
	}
	testDBOnce.Do(func() { testDB, testDBErr = testdb.Open(url, testDBCfg) })
	if testDBErr != nil {
		t.Fatalf("test database: %v", testDBErr)
	}
	if err := testdb.Truncate(testDB, testDBCfg); err != nil {
		t.Fatalf("reset test database: %v", err)
	}
	return testDB
}

func newTestEraser(t *testing.T, db *gorm.DB, client *matrix.Client) *Eraser {
	t.Helper()
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	eraser, err := New(context.Background(), Deps{
		DB: db, SQL: sqlDB, Matrix: client, Issuer: testIssuer, Bot: testBot, Logger: quiet,
	})
	if err != nil {
		t.Fatalf("new eraser: %v", err)
	}
	return eraser
}

func mustExec(t *testing.T, db *gorm.DB, query string, args ...any) {
	t.Helper()
	if err := db.Exec(query, args...).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
}

// seedPerson creates an account with a row in every table that references a
// person: identity, session, push device, Matrix account and session, agent,
// a space they created and a call they started.
func seedPerson(t *testing.T, db *gorm.DB, local string) person {
	t.Helper()
	p := person{email: local + "@facile.studio", subject: "sub-" + local,
		mxid: "@" + local + ":" + testServer, ag: "@" + local + "-bot:" + testServer}
	if err := db.Raw(`INSERT INTO porte_users (email, name, avatar_url) VALUES (?, ?, 'https://sso.test/avatars/x')
		RETURNING id`, p.email, local).Scan(&p.id).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	mustExec(t, db, `INSERT INTO porte_identities (user_id, provider, subject) VALUES (?, ?, ?)`, p.id, testIssuer, p.subject)
	var session int64
	mustExec(t, db, `INSERT INTO porte_sessions (token_hash, user_id) VALUES (?, ?)`, "hash-"+local, p.id)
	db.Raw(`SELECT id FROM porte_sessions WHERE user_id = ?`, p.id).Scan(&session)
	mustExec(t, db, `INSERT INTO push_devices (token, user_id, porte_session_id) VALUES (?, ?, ?)`, "fcm-"+local, p.id, session)
	rows := []any{
		&schemas.Account{UserID: p.id, Localpart: local, MatrixUserID: p.mxid, AvatarURL: "mxc://hs/" + local},
		&schemas.MatrixSession{PorteSessionID: session, UserID: p.id, DeviceID: "D" + local, SealedToken: "sealed"},
		&schemas.Agent{OwnerID: p.id, Slug: local + "-bot", Name: "Bot of " + local, MatrixUserID: p.ag},
		&schemas.Space{RoomID: "!space-" + local + ":hs", Name: "Space " + local, Slug: "space-" + local, CreatedBy: p.id},
		&schemas.Call{RoomID: "!call-" + local + ":hs", CallID: "c", EventID: "e", StartedBy: p.mxid},
	}
	for _, row := range rows {
		if err := db.Create(row).Error; err != nil {
			t.Fatalf("seed %T: %v", row, err)
		}
	}
	return p
}
