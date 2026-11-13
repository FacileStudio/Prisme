package schemas

import (
	"errors"
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

var errRollBack = errors.New("roll back")

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
		FROM information_schema.columns WHERE table_schema = current_schema() AND table_name = @table
		UNION ALL SELECT replace(indexdef, ' ON ' || schemaname || '.', ' ON ')
		FROM pg_indexes WHERE schemaname = current_schema() AND tablename = @table
		UNION ALL SELECT conname || ' ' || pg_get_constraintdef(oid) FROM pg_constraint
		WHERE conrelid = to_regclass(@table)
		ORDER BY 1`, map[string]any{"table": table}).Scan(&shape).Error
	if err != nil {
		t.Fatalf("read the shape of %s: %v", table, err)
	}
	return shape
}

// assertMatchesLibrary builds a library's own DDL in a scratch schema, inside
// a transaction that is rolled back, and compares each table's columns,
// types, defaults, indexes and constraints with the ones the migrations
// built. Running the DDL over the migrated schema instead proves nothing,
// since its CREATE TABLE IF NOT EXISTS skips a table that already exists.
func assertMatchesLibrary(t *testing.T, db *gorm.DB, ddl string, tables []string) {
	t.Helper()
	migrated := map[string][]string{}
	for _, table := range tables {
		migrated[table] = tableShape(t, db, table)
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, statement := range []string{`CREATE SCHEMA prisme_library_drift`, `SET LOCAL search_path TO prisme_library_drift`, ddl} {
			if err := tx.Exec(statement).Error; err != nil {
				t.Fatalf("build the library schema: %v", err)
			}
		}
		for _, table := range tables {
			want := tableShape(t, tx, table)
			if len(want) == 0 {
				t.Errorf("the library no longer creates %s", table)
			} else if !slices.Equal(migrated[table], want) {
				t.Errorf("%s drifted from the library:\nmigrations %v\nlibrary    %v", table, migrated[table], want)
			}
		}
		return errRollBack
	})
	if !errors.Is(err, errRollBack) {
		t.Fatalf("transaction: %v", err)
	}
}

// TestPorteSchemaIsInTheMigrations checks the migrations build porte's tables
// exactly as porte/pg's Schema does, which Prisme used to execute at every
// boot. A porte bump that changes it makes this fail, and the change then
// needs a numbered migration of its own.
func TestPorteSchemaIsInTheMigrations(t *testing.T) {
	assertMatchesLibrary(t, openSchemasDatabase(t), pg.Schema, porteTables)
}
