package users

import (
	"context"
	"slices"
	"testing"

	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/testdb"

	"gorm.io/gorm"
)

const usersBootstrap = "docker run -d -e POSTGRES_PASSWORD=t -e POSTGRES_DB=t -p 5432:5432 postgres:16"

var usersDBCfg = testdb.Config{Prefix: "prisme_users_test", Migrate: schemas.Migrate}

func openUsersDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	url, configured := testdb.URL()
	if !configured {
		testdb.Announce(usersBootstrap)
		t.Skip(testdb.SkipReason(usersBootstrap))
	}
	db, err := testdb.Open(url, usersDBCfg)
	if err != nil {
		t.Fatalf("test database: %v", err)
	}
	if err := testdb.Truncate(db, usersDBCfg); err != nil {
		t.Fatalf("reset test database: %v", err)
	}
	return db
}

type fakeRooms map[string][]string

func (f fakeRooms) JoinedRooms(_ context.Context, userID string) ([]string, error) {
	var rooms []string
	for room, members := range f {
		if slices.Contains(members, userID) {
			rooms = append(rooms, room)
		}
	}
	return rooms, nil
}

func (f fakeRooms) JoinedMembers(_ context.Context, _, roomID string) ([]string, error) {
	return f[roomID], nil
}

func seedPerson(t *testing.T, db *gorm.DB, id int64, localpart string) string {
	t.Helper()
	mxid := "@" + localpart + ":prisme.test"
	if err := db.Exec(`INSERT INTO porte_users (id, email, name) VALUES (?, ?, ?)`,
		id, localpart+"@example.com", localpart).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&schemas.Account{UserID: id, Localpart: localpart, MatrixUserID: mxid}).Error; err != nil {
		t.Fatal(err)
	}
	return mxid
}

func seedAgent(t *testing.T, db *gorm.DB, owner int64, slug string) string {
	t.Helper()
	mxid := "@" + slug + ":prisme.test"
	if err := db.Create(&schemas.Agent{OwnerID: owner, Slug: slug, Name: slug, MatrixUserID: mxid}).Error; err != nil {
		t.Fatal(err)
	}
	return mxid
}

func seedSpace(t *testing.T, db *gorm.DB, roomID, slug string) {
	t.Helper()
	if err := db.Create(&schemas.Space{RoomID: roomID, Name: slug, Slug: slug, CreatedBy: 1}).Error; err != nil {
		t.Fatal(err)
	}
}

func userIDs(users []User) []string {
	ids := make([]string, 0, len(users))
	for _, user := range users {
		ids = append(ids, user.UserID)
	}
	slices.Sort(ids)
	return ids
}

// directory is two spaces: alice shares "team" with bob and his agent robo;
// carol and her agent spy sit in "other"; dave shares only a DM with alice.
// alice's agent pixel joined nothing.
type directory struct {
	service                          *Service
	alice, bob, aliceAgent, bobAgent string
}

func seedDirectory(t *testing.T) directory {
	t.Helper()
	db := openUsersDatabase(t)
	d := directory{alice: seedPerson(t, db, 1, "alice"), bob: seedPerson(t, db, 2, "bob")}
	carol, dave := seedPerson(t, db, 3, "carol"), seedPerson(t, db, 4, "dave")
	d.aliceAgent, d.bobAgent = seedAgent(t, db, 1, "pixel"), seedAgent(t, db, 2, "robo")
	carolAgent := seedAgent(t, db, 3, "spy")
	seedSpace(t, db, "!team:prisme.test", "team")
	seedSpace(t, db, "!other:prisme.test", "other")
	d.service = NewService(db, fakeRooms{
		"!team:prisme.test":  {d.alice, d.bob, d.bobAgent},
		"!other:prisme.test": {carol, carolAgent},
		"!dm:prisme.test":    {d.alice, dave},
	})
	return d
}

func TestDirectoryShowsOnlySpacePeers(t *testing.T) {
	d := seedDirectory(t)
	found, err := d.service.Search(context.Background(), 1, "")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{d.aliceAgent, d.alice, d.bob, d.bobAgent}
	slices.Sort(want)
	if got := userIDs(found); !slices.Equal(got, want) {
		t.Fatalf("alice sees %v, want %v", got, want)
	}
	if found, err := d.service.Search(context.Background(), 99, ""); err != nil || len(found) != 0 {
		t.Fatalf("a caller with no account found %v, %v; want nobody", found, err)
	}
}

func TestDirectoryFindsNobodyOutsideTheCallersSpaces(t *testing.T) {
	d := seedDirectory(t)
	for _, query := range []string{"carol", "carol@example.com", "dave", "spy"} {
		found, err := d.service.Search(context.Background(), 1, query)
		if err != nil {
			t.Fatal(err)
		}
		if len(found) != 0 {
			t.Fatalf("search %q as alice found %v, want nobody outside her spaces", query, userIDs(found))
		}
	}
}

func TestAgentListShowsOnlyOwnAndPeerAgents(t *testing.T) {
	d := seedDirectory(t)
	agents, err := d.service.Agents(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	var seen []string
	for _, agent := range agents {
		seen = append(seen, agent.UserID+" "+agent.OwnerUserID)
	}
	slices.Sort(seen)
	if want := []string{d.aliceAgent + " " + d.alice, d.bobAgent + " " + d.bob}; !slices.Equal(seen, want) {
		t.Fatalf("alice sees agents %v, want %v", seen, want)
	}
}
