package erasure

import (
	"testing"

	"gorm.io/gorm"
)

const textScan = `
SELECT table_name || '.' || column_name AS place, n FROM (
  SELECT c.table_name, c.column_name,
    (xpath('/row/n/text()', query_to_xml(format(
      'SELECT count(*) AS n FROM %I.%I WHERE %I::text IN (%L, %L, %L) OR lower(%I::text) = lower(%L)',
      c.table_schema, c.table_name, c.column_name, ?::text, ?::text, ?::text, c.column_name, ?::text),
      false, true, '')))[1]::text::int AS n
  FROM information_schema.columns c
  JOIN information_schema.tables t USING (table_schema, table_name)
  WHERE t.table_type = 'BASE TABLE' AND c.table_schema = current_schema()
    AND c.data_type IN ('text', 'character varying', 'uuid')
) counts WHERE n > 0`

const idScan = `
SELECT table_name || '.' || column_name AS place, n FROM (
  SELECT c.table_name, c.column_name,
    (xpath('/row/n/text()', query_to_xml(format('SELECT count(*) AS n FROM %I.%I WHERE %I = %s',
      c.table_schema, c.table_name, c.column_name, ?::bigint), false, true, '')))[1]::text::int AS n
  FROM information_schema.columns c
  WHERE c.table_schema = current_schema() AND c.data_type = 'bigint'
    AND (c.column_name IN ('user_id', 'owner_id', 'created_by')
      OR (c.table_name = 'porte_users' AND c.column_name = 'id'))
) counts WHERE n > 0`

type leftover struct {
	Place string
	N     int64
}

// leftovers lists every column of the test schema that still holds the
// person's email, subject, Matrix ids or account id, with a count.
func leftovers(t *testing.T, db *gorm.DB, p person) []leftover {
	t.Helper()
	var text, ids []leftover
	if err := db.Raw(textScan, p.subject, p.mxid, p.ag, p.email).Scan(&text).Error; err != nil {
		t.Fatalf("scan text columns: %v", err)
	}
	if err := db.Raw(idScan, p.id).Scan(&ids).Error; err != nil {
		t.Fatalf("scan id columns: %v", err)
	}
	return append(text, ids...)
}

func assertNoLeftovers(t *testing.T, db *gorm.DB, p person) {
	t.Helper()
	if found := leftovers(t, db, p); len(found) > 0 {
		t.Errorf("rows still reference the erased person: %+v", found)
	}
}

func assertUntouched(t *testing.T, db *gorm.DB, p person) {
	t.Helper()
	found := map[string]int64{}
	for _, l := range leftovers(t, db, p) {
		found[l.Place] = l.N
	}
	for _, place := range []string{"porte_users.email", "porte_identities.subject", "accounts.matrix_user_id",
		"agents.matrix_user_id", "matrix_sessions.user_id", "push_devices.user_id", "spaces.created_by"} {
		if found[place] != 1 {
			t.Errorf("%s holds %d rows of the other person, want 1", place, found[place])
		}
	}
}
