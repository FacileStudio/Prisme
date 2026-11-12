package erasure

import (
	"context"
	"io"
	"log/slog"
	"slices"
	"testing"

	"github.com/FacileStudio/Prisme/apps/api/modules/spaces"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/pg"
)

// recreate seeds a second account for the erased person's subject under
// the issuer, as a service call racing the erasure made one before v0.14.2.
func recreate(t *testing.T, w world, issuer string) person {
	t.Helper()
	back := seedPerson(t, w.db, "back")
	back.subject = w.gone.subject
	mustExec(t, w.db, `UPDATE porte_identities SET provider = ?, subject = ? WHERE user_id = ?`, issuer, back.subject, back.id)
	w.synapse.addUser(back.mxid, 1)
	w.synapse.addUser(back.ag, 0)
	return back
}

func assertRecreatedErased(t *testing.T, w world, back person) {
	t.Helper()
	assertNoLeftovers(t, w.db, back)
	assertUntouched(t, w.db, w.other)
	w.synapse.read(func(f *fakeSynapse) {
		for _, id := range []string{back.mxid, back.ag} {
			if !slices.Contains(f.calls, "POST /v1/deactivate/"+id) || !f.users[id] {
				t.Errorf("Synapse did not deactivate %s: %v", id, f.calls)
			}
		}
	})
}

// A notice handled before an account was re-created for its subject is
// handled again on the next delivery to the same running receiver: porte asks
// Holds, which finds the new identity, and erases it instead of skipping the
// subject as done.
func TestARepushedNoticeErasesAnAccountRecreatedAfterTheErasure(t *testing.T) {
	w := seedWorld(t)
	issuer := fakeRegistre(t, w.gone.subject)
	mustExec(t, w.db, `UPDATE porte_identities SET provider = ?`, issuer)
	sqlDB, _ := w.db.DB()
	store := pg.New(sqlDB)
	receiver := &directory.Deletions{
		Client:          directory.Client{Issuer: issuer, ClientID: "prisme", ClientSecret: "s"},
		OnPersonDeleted: w.eraser.ErasePerson, Erase: store.EraseIdentity, Holds: store.HoldsIdentity,
		State: spaces.NewDeletionState(w.db), Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	if err := receiver.Poll(context.Background()); err != nil {
		t.Fatalf("first notice: %v", err)
	}
	assertNoLeftovers(t, w.db, w.gone)
	back := recreate(t, w, issuer)
	if err := receiver.Poll(context.Background()); err != nil {
		t.Fatalf("repushed notice: %v", err)
	}
	assertRecreatedErased(t, w, back)
}

// erase-user takes an account whose only subject Prisme already erased: no
// notice reaches it, so the operator's path must.
func TestEraseUserErasesAnAccountRecreatedAfterTheErasure(t *testing.T) {
	w := seedWorld(t)
	if err := w.erase(t); err != nil {
		t.Fatalf("erase: %v", err)
	}
	back := recreate(t, w, testIssuer)
	if found, err := w.eraser.EraseUnlinked(context.Background(), back.email); !found || err != nil {
		t.Fatalf("EraseUnlinked = %v, %v", found, err)
	}
	assertRecreatedErased(t, w, back)
}
