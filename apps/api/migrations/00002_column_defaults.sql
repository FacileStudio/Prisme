-- +goose Up
-- GORM filled these timestamps itself, so AutoMigrate never gave them a database default and a
-- row written outside GORM got NULL. The database now fills them.

ALTER TABLE accounts ALTER COLUMN created_at SET DEFAULT now();
ALTER TABLE agents ALTER COLUMN created_at SET DEFAULT now();
ALTER TABLE calls ALTER COLUMN started_at SET DEFAULT now();
ALTER TABLE channels ALTER COLUMN created_at SET DEFAULT now();
ALTER TABLE erased_subjects ALTER COLUMN erased_at SET DEFAULT now();
ALTER TABLE matrix_sessions ALTER COLUMN created_at SET DEFAULT now();
ALTER TABLE spaces ALTER COLUMN created_at SET DEFAULT now();

-- +goose Down
ALTER TABLE accounts ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE agents ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE calls ALTER COLUMN started_at DROP DEFAULT;
ALTER TABLE channels ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE erased_subjects ALTER COLUMN erased_at DROP DEFAULT;
ALTER TABLE matrix_sessions ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE spaces ALTER COLUMN created_at DROP DEFAULT;
