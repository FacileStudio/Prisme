package spaces

import (
	"sync"
	"testing"

	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/testdb"

	"gorm.io/gorm"
)

const testBootstrap = "docker run -d -e POSTGRES_PASSWORD=t -e POSTGRES_DB=t -p 5432:5432 postgres:16"

var (
	testDBOnce sync.Once
	testDB     *gorm.DB
	testDBErr  error
	testDBCfg  = testdb.Config{Prefix: "prisme_spaces_test", Migrate: schemas.Migrate}
)

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
