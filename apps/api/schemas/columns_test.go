package schemas

import (
	"reflect"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/tronc/testdb"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const schemasBootstrap = "docker run -d -e POSTGRES_PASSWORD=t -e POSTGRES_DB=t -p 5432:5432 postgres:16"

var schemasDBCfg = testdb.Config{Prefix: "prisme_schemas_test", Migrate: Migrate}

var queryModels = []any{
	&Account{}, &MatrixSession{}, &Space{}, &Channel{}, &Agent{}, &Call{}, &ErasedSubject{},
}

var porteTables = []string{
	"porte_users", "porte_identities", "porte_sessions", "porte_login_codes",
	"porte_magic_links", "porte_api_keys", "porte_api_key_usage",
}

var timeType = reflect.TypeFor[time.Time]()

func openSchemasDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	url, configured := testdb.URL()
	if !configured {
		testdb.Announce(schemasBootstrap)
		t.Skip(testdb.SkipReason(schemasBootstrap))
	}
	db, err := testdb.Open(url, schemasDBCfg)
	if err != nil {
		t.Fatalf("test database: %v", err)
	}
	return db
}

func compatibleTypes(field reflect.Type) []string {
	for field.Kind() == reflect.Pointer {
		field = field.Elem()
	}
	if field == timeType {
		return []string{"timestamp with time zone"}
	}
	switch field.Kind() {
	case reflect.String:
		return []string{"text", "character varying"}
	case reflect.Bool:
		return []string{"boolean"}
	case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
		return []string{"bigint"}
	case reflect.Int32, reflect.Uint32:
		return []string{"integer", "bigint"}
	default:
		return nil
	}
}

func migratedColumnTypes(t *testing.T, db *gorm.DB, table string) map[string]string {
	t.Helper()
	var rows []struct {
		ColumnName string
		DataType   string
	}
	err := db.Raw(`SELECT column_name, data_type FROM information_schema.columns
		WHERE table_schema = current_schema() AND table_name = ?`, table).Scan(&rows).Error
	if err != nil {
		t.Fatalf("read the columns of %s: %v", table, err)
	}
	columns := map[string]string{}
	for _, row := range rows {
		columns[row.ColumnName] = row.DataType
	}
	return columns
}

// TestGORMModelsMatchTheMigratedSchema checks every column a GORM model reads
// or writes exists in the schema the migrations build, with a type the field
// can hold. AutoMigrate no longer adds a missing column, so a field added to a
// struct without a migration fails here instead of at query time.
func TestGORMModelsMatchTheMigratedSchema(t *testing.T) {
	db := openSchemasDatabase(t)
	for _, model := range queryModels {
		parsed, err := schema.Parse(model, &sync.Map{}, db.NamingStrategy)
		if err != nil {
			t.Fatalf("parse %T: %v", model, err)
		}
		columns := migratedColumnTypes(t, db, parsed.Table)
		if len(columns) == 0 {
			t.Errorf("%T: table %s does not exist", model, parsed.Table)
			continue
		}
		for _, field := range parsed.Fields {
			if field.DBName == "" {
				continue
			}
			got, ok := columns[field.DBName]
			if !ok {
				t.Errorf("%s.%s: no such column", parsed.Table, field.DBName)
				continue
			}
			if want := compatibleTypes(field.FieldType); !slices.Contains(want, got) {
				t.Errorf("%s.%s: column is %s, field %s needs one of %v", parsed.Table, field.DBName, got, field.FieldType, want)
			}
		}
	}
}

func tableShape(t *testing.T, db *gorm.DB, table string) []string {
	t.Helper()
	var shape []string
	err := db.Raw(`SELECT column_name || ' ' || data_type || ' ' || is_nullable || ' ' || coalesce(column_default, '')
		FROM information_schema.columns WHERE table_schema = current_schema() AND table_name = ?
		UNION ALL SELECT indexdef FROM pg_indexes WHERE schemaname = current_schema() AND tablename = ?
		UNION ALL SELECT conname || ' ' || pg_get_constraintdef(oid) FROM pg_constraint
		 WHERE conrelid = to_regclass(?)
		ORDER BY 1`, table, table, table).Scan(&shape).Error
	if err != nil {
		t.Fatalf("read the shape of %s: %v", table, err)
	}
	return shape
}

// TestPorteSchemaIsInTheMigrations checks the migrations already hold
// porte/pg's Schema, which Prisme used to execute at every boot. A porte bump
// that changes it makes this fail, and the change then needs a numbered
// migration of its own.
func TestPorteSchemaIsInTheMigrations(t *testing.T) {
	db := openSchemasDatabase(t)
	before := map[string][]string{}
	for _, table := range porteTables {
		before[table] = tableShape(t, db, table)
		if len(before[table]) == 0 {
			t.Fatalf("the migrations do not create %s", table)
		}
	}
	if err := db.Exec(pg.Schema).Error; err != nil {
		t.Fatalf("run porte's Schema: %v", err)
	}
	for _, table := range porteTables {
		if after := tableShape(t, db, table); !slices.Equal(before[table], after) {
			t.Errorf("porte's Schema changed %s:\nbefore %v\nafter  %v", table, before[table], after)
		}
	}
}
