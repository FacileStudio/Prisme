package erasure

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"slices"
	"strings"
	"testing"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"

	"gorm.io/gorm"
)

func seedLive(t *testing.T, db *gorm.DB, gone *person) {
	t.Helper()
	if err := db.Raw(`INSERT INTO porte_users (email) VALUES (?) RETURNING id`, gone.email).Scan(&gone.id).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	mustExec(t, db, `INSERT INTO porte_identities (user_id, provider, subject) VALUES (?, ?, ?)`, gone.id, testIssuer, gone.subject)
	mustExec(t, db, `INSERT INTO accounts (user_id, localpart, matrix_user_id, avatar_url) VALUES (?, ?, ?, '')`,
		gone.id, localOf(gone.mxid), gone.mxid)
	mustExec(t, db, `INSERT INTO agents (owner_id, slug, name, matrix_user_id) VALUES (?, ?, 'Bot', ?)`,
		gone.id, localOf(gone.ag), gone.ag)
}

func localOf(id string) string {
	local, _, _ := strings.Cut(strings.TrimPrefix(id, "@"), ":")
	return local
}

// Against a real Synapse: the person and their agent are deactivated and
// erased, the room only they were in is purged, the direct chat they alone
// administered keeps an admin, the colleague's room stays, the profile is
// gone, the id is never handed out again and a newcomer with the same email
// gets a new Matrix account.
func TestLiveErasureAgainstSynapse(t *testing.T) {
	client := liveSynapse(t)
	db := openTestDatabase(t)
	gone, colleague, rooms := liveWorld(t, client, liveRun())
	seedLive(t, db, &gone)
	eraser := newTestEraser(t, db, client)
	eraser.deps.Bot = client.UserID("prisme")
	mustLive(t, eraser.ErasePerson(context.Background(), testIssuer, gone.subject), "erase")
	assertLiveErased(t, client, gone, colleague, rooms)
	assertNoLeftovers(t, db, gone)
	assertLiveNewcomer(t, db, client, gone)
}

func assertLiveErased(t *testing.T, client *matrix.Client, gone person, colleague string, rooms liveRooms) {
	t.Helper()
	ctx := context.Background()
	for _, id := range []string{gone.mxid, gone.ag} {
		if retired, err := client.Retired(ctx, id); err != nil || !retired {
			t.Errorf("%s retired = %v, %v", id, retired, err)
		}
	}
	if _, err := client.RoomMembers(ctx, rooms.solo); !matrix.IsCode(err, "M_NOT_FOUND") {
		t.Errorf("the solo room was not purged: %v", err)
	}
	if level := powerOf(t, client, rooms.direct, colleague); level != 100 {
		t.Errorf("the colleague holds %d in the direct chat, want 100", level)
	}
	if members := waitParted(t, client, rooms.shared, gone.mxid); !slices.Equal(members, []string{colleague}) {
		t.Errorf("shared room members = %v, want only the colleague", members)
	}
	if name := client.DisplayName(ctx, gone.mxid); name != gone.mxid {
		t.Errorf("the erased profile still answers %q", name)
	}
}

func assertLiveNewcomer(t *testing.T, db *gorm.DB, client *matrix.Client, gone person) {
	t.Helper()
	local := localOf(gone.mxid)
	if err := client.Register(context.Background(), local); !errors.Is(err, matrix.ErrRetired) {
		t.Fatalf("registering the erased id = %v, want ErrRetired", err)
	}
	var newcomer int64
	mustLive(t, db.Raw(`INSERT INTO porte_users (email) VALUES (?) RETURNING id`, gone.email).Scan(&newcomer).Error, "newcomer")
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	account, err := accounts.NewService(db, client, quiet).Ensure(context.Background(), newcomer)
	mustLive(t, err, "ensure the newcomer")
	if account.MatrixUserID == gone.mxid || account.Localpart != local+"2" {
		t.Errorf("the newcomer got %s, want %s2", account.MatrixUserID, local)
	}
}
