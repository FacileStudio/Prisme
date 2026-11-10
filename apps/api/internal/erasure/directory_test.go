package erasure

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/pg"
)

type memoryState struct {
	mu    sync.Mutex
	saved []byte
}

func (m *memoryState) Load(context.Context) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.saved, nil
}

func (m *memoryState) Save(_ context.Context, state []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.saved = state
	return nil
}

// fakeRegistre lists the given subjects as erased a minute ago.
func fakeRegistre(t *testing.T, subjects ...string) string {
	t.Helper()
	at := time.Now().Add(-time.Minute).UTC()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		list := make([]directory.Deleted, 0, len(subjects))
		for _, subject := range subjects {
			list = append(list, directory.Deleted{Subject: subject, DeletedAt: at})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"deleted": list})
	}))
	t.Cleanup(server.Close)
	return server.URL
}

func receiver(t *testing.T, w world, issuer string) *directory.Deletions {
	t.Helper()
	sqlDB, _ := w.db.DB()
	return &directory.Deletions{
		Client:          directory.Client{Issuer: issuer, ClientID: "prisme", ClientSecret: "s"},
		OnPersonDeleted: w.eraser.ErasePerson, Erase: pg.New(sqlDB).EraseIdentity,
		State: &memoryState{}, Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

// The poll the directory syncer runs reaches the same erasure: the person is
// gone, the notice is recorded as handled and a second poll erases nothing.
func TestTheDeletionPollErasesThePerson(t *testing.T) {
	w := seedWorld(t)
	issuer := fakeRegistre(t, w.gone.subject)
	mustExec(t, w.db, `UPDATE porte_identities SET provider = ?`, issuer)
	deletions := receiver(t, w, issuer)
	if err := deletions.Poll(context.Background()); err != nil {
		t.Fatalf("poll: %v", err)
	}
	assertNoLeftovers(t, w.db, w.gone)
	var before int
	w.synapse.read(func(f *fakeSynapse) { before = len(f.calls) })
	if err := deletions.Poll(context.Background()); err != nil {
		t.Fatalf("second poll: %v", err)
	}
	w.synapse.read(func(f *fakeSynapse) {
		if len(f.calls) != before {
			t.Errorf("the second poll called Synapse again: %v", f.calls[before:])
		}
	})
}

// A notice that cannot be proven done is not recorded, so Registre's list
// keeps offering it to the next poll.
func TestAnUnprovenNoticeStaysPending(t *testing.T) {
	w := seedWorld(t)
	issuer := fakeRegistre(t, "never-here")
	mustExec(t, w.db, `UPDATE porte_identities SET provider = ?`, issuer)
	mustExec(t, w.db, `INSERT INTO porte_users (email) VALUES ('service-only@facile.studio')`)
	deletions := receiver(t, w, issuer)
	for range 2 {
		if err := deletions.Poll(context.Background()); err == nil {
			t.Fatal("a poll answered an unproven notice as done")
		}
	}
}

// Matrix never reissues an id. A new person who signs in with the email of an
// erased one gets a new, unrelated Matrix account, never the erased one.
func TestANewPersonWithAnErasedEmailGetsANewMatrixAccount(t *testing.T) {
	w := seedWorld(t)
	if err := w.erase(t); err != nil {
		t.Fatalf("erase: %v", err)
	}
	var newcomer int64
	if err := w.db.Raw(`INSERT INTO porte_users (email) VALUES (?) RETURNING id`, w.gone.email).Scan(&newcomer).Error; err != nil {
		t.Fatalf("new user: %v", err)
	}
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	account, err := accounts.NewService(w.db, w.eraser.deps.Matrix, quiet).Ensure(context.Background(), newcomer)
	if err != nil {
		t.Fatalf("ensure: %v", err)
	}
	if account.MatrixUserID == w.gone.mxid || account.Localpart != "gone2" {
		t.Fatalf("the newcomer got %s, want a new id beside the erased %s", account.MatrixUserID, w.gone.mxid)
	}
	w.synapse.read(func(f *fakeSynapse) {
		if deactivated, ok := f.users[account.MatrixUserID]; !ok || deactivated {
			t.Errorf("the new account %s is not an active Synapse user", account.MatrixUserID)
		}
	})
}
