package erasure

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"testing"

	"github.com/FacileStudio/porte"

	"gorm.io/gorm"
)

// world is a person to erase, another person who shares rooms with them, the
// fake homeserver and the eraser.
type world struct {
	db          *gorm.DB
	synapse     *fakeSynapse
	eraser      *Eraser
	gone, other person
}

func seedWorld(t *testing.T) world {
	t.Helper()
	db := openTestDatabase(t)
	synapse, client := newFakeSynapse(t)
	w := world{db: db, synapse: synapse, eraser: newTestEraser(t, db, client),
		gone: seedPerson(t, db, "gone"), other: seedPerson(t, db, "other")}
	synapse.addUser(w.gone.mxid, 2)
	synapse.addUser(w.gone.ag, 1)
	synapse.addUser(w.other.mxid, 1)
	synapse.addUser(testBot, 0)
	synapse.addRoom("!solo:hs", w.gone.mxid, "12", nil, w.gone.mxid, w.gone.ag)
	synapse.addRoom("!dm:hs", w.gone.mxid, "12", nil, w.gone.mxid, w.other.mxid)
	synapse.addRoom("!shared:hs", w.other.mxid, "12", nil, w.gone.mxid, w.other.mxid)
	synapse.addRoom("!space-gone:hs", w.gone.mxid, "12", map[string]int64{testBot: 100}, w.gone.mxid, testBot)
	return w
}

func (w world) erase(t *testing.T) error {
	t.Helper()
	return w.eraser.ErasePerson(context.Background(), testIssuer, w.gone.subject)
}

func countRows(t *testing.T, db *gorm.DB, query string, args ...any) int64 {
	t.Helper()
	var n int64
	if err := db.Raw(query, args...).Scan(&n).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

// The person's own rows go, every Matrix account they own is erased with its
// media, the room nobody else is in is purged, the shared room they alone
// administered gains an admin, and the shared space and the call stay with no
// trace of them. The other person keeps everything.
func TestErasureRemovesThePersonFromPrismeAndSynapse(t *testing.T) {
	w := seedWorld(t)
	if err := w.erase(t); err != nil {
		t.Fatalf("erase: %v", err)
	}
	assertNoLeftovers(t, w.db, w.gone)
	assertUntouched(t, w.db, w.other)
	w.synapse.read(func(f *fakeSynapse) { assertSynapseErased(t, f, w) })
	if n := countRows(t, w.db, `SELECT count(*) FROM spaces WHERE room_id = '!space-gone:hs' AND created_by = 0`); n != 1 {
		t.Error("the shared space lost its row or kept its creator")
	}
	if n := countRows(t, w.db, `SELECT count(*) FROM calls WHERE room_id = '!call-gone:hs' AND started_by = ''`); n != 1 {
		t.Error("the running call lost its row or kept its starter")
	}
}

// Registre delivers at least once. A notice for a person already erased is
// done, even while an account with no subject would otherwise block it.
func TestARepeatedNoticeIsAnsweredAsDone(t *testing.T) {
	w := seedWorld(t)
	if err := w.erase(t); err != nil {
		t.Fatalf("erase: %v", err)
	}
	mustExec(t, w.db, `INSERT INTO porte_users (email) VALUES ('service-only@facile.studio')`)
	if err := w.erase(t); err != nil {
		t.Fatalf("second notice = %v, want nil", err)
	}
}

// A subject no row holds is done only when every account holds a subject. An
// account a service key created by email might be that person, so the notice
// stays pending until the operator erases or the person links it.
func TestAnUnknownSubjectIsDoneOnlyWhenEveryAccountIsLinked(t *testing.T) {
	w := seedWorld(t)
	ctx := context.Background()
	if err := w.eraser.ErasePerson(ctx, testIssuer, "never-here"); err != nil {
		t.Fatalf("unknown subject with every account linked = %v, want nil", err)
	}
	mustExec(t, w.db, `INSERT INTO porte_users (email) VALUES ('Service-Only@facile.studio')`)
	if err := w.eraser.ErasePerson(ctx, testIssuer, "never-here"); !errors.Is(err, ErrUnproven) {
		t.Fatalf("unknown subject beside an unlinked account = %v, want ErrUnproven", err)
	}
	if _, err := w.eraser.EraseUnlinked(ctx, w.gone.email); err == nil {
		t.Fatal("EraseUnlinked erased a linked account")
	}
	if found, err := w.eraser.EraseUnlinked(ctx, "service-only@facile.studio"); !found || err != nil {
		t.Fatalf("EraseUnlinked = %v, %v", found, err)
	}
	if err := w.eraser.ErasePerson(ctx, testIssuer, "never-here"); err != nil {
		t.Fatalf("unknown subject after the operator erased the unlinked account = %v", err)
	}
}

// Synapse failing first leaves every row, so the retry finds the person again.
func TestASynapseFailureKeepsTheRowsForTheRetry(t *testing.T) {
	w := seedWorld(t)
	w.synapse.failErase(true)
	if err := w.erase(t); err == nil {
		t.Fatal("erase succeeded while Synapse refused the deactivation")
	}
	if len(leftovers(t, w.db, w.gone)) == 0 {
		t.Fatal("rows were deleted although Synapse kept the account")
	}
	w.synapse.failErase(false)
	if err := w.erase(t); err != nil {
		t.Fatalf("retry: %v", err)
	}
	assertNoLeftovers(t, w.db, w.gone)
}

// The database failing after Synapse erased the accounts leaves the rows, and
// the retry repeats the Synapse calls harmlessly and finishes the rows.
func TestADatabaseFailureAfterSynapseConvergesOnRetry(t *testing.T) {
	w := seedWorld(t)
	w.eraser.store.EraseData = func(context.Context, *sql.Tx, int64) error { return errors.New("disk full") }
	if err := w.erase(t); err == nil {
		t.Fatal("erase succeeded while the transaction failed")
	}
	var deactivated bool
	w.synapse.read(func(f *fakeSynapse) { deactivated = f.users[w.gone.mxid] })
	if !deactivated || len(leftovers(t, w.db, w.gone)) == 0 {
		t.Fatal("want Synapse erased and the rows kept after a failed transaction")
	}
	w.eraser.store.EraseData = eraseRows
	if err := w.erase(t); err != nil {
		t.Fatalf("retry: %v", err)
	}
	assertNoLeftovers(t, w.db, w.gone)
}

func assertSynapseErased(t *testing.T, f *fakeSynapse, w world) {
	t.Helper()
	for _, id := range []string{w.gone.mxid, w.gone.ag} {
		if !f.users[id] || f.media[id] != 0 {
			t.Errorf("%s: deactivated %v, media left %d", id, f.users[id], f.media[id])
		}
	}
	if f.users[w.other.mxid] {
		t.Error("the other person was deactivated")
	}
	if _, kept := f.rooms["!solo:hs"]; kept {
		t.Error("the room only the person and their agent were in was not purged")
	}
	if f.rooms["!dm:hs"].levels[w.other.mxid] != 100 {
		t.Error("the shared room the person alone administered was left without an admin")
	}
	for _, kept := range []string{"!shared:hs", "!space-gone:hs"} {
		if _, ok := f.rooms[kept]; !ok || slices.Contains(f.calls, "DELETE /v1/rooms/"+kept) {
			t.Errorf("%s was purged; it belongs to its other members", kept)
		}
	}
}

// The erase store sits on porte_users, which porte/pg's UserStore owns. Its
// insert would leave only facile_id NULL, and it refuses to insert at all.
func TestTheEraseStoreNeverCreatesAnAccount(t *testing.T) {
	db := openTestDatabase(t)
	_, client := newFakeSynapse(t)
	eraser := newTestEraser(t, db, client)
	if got := eraser.store.Unwritten(); !slices.Equal(got, []string{"facile_id"}) {
		t.Errorf("Unwritten = %v, want [facile_id]", got)
	}
	claims := porte.Claims{Provider: testIssuer, Subject: "new", Email: "new@facile.studio", EmailVerified: true}
	if _, err := eraser.store.UpsertFromOIDC(context.Background(), claims); err == nil {
		t.Fatal("the erase store created an account")
	}
	if n := countRows(t, db, `SELECT count(*) FROM porte_users`); n != 0 {
		t.Fatalf("porte_users holds %d rows after a refused creation", n)
	}
}
