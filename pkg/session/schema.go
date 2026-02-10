package session

import "database/sql"

const schemaVersion = 1

const ddl = `
CREATE TABLE IF NOT EXISTS schema_version (
    version INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS buffer (
    row_index INTEGER PRIMARY KEY,
    text      TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS cursor (
    id  INTEGER PRIMARY KEY CHECK (id = 1),
    row INTEGER NOT NULL DEFAULT 0,
    col INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS variables (
    name      TEXT PRIMARY KEY,
    kind      INTEGER NOT NULL,
    num       REAL    NOT NULL DEFAULT 0,
    str       TEXT    NOT NULL DEFAULT '',
    type_code TEXT    NOT NULL DEFAULT '',
    period    INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS results (
    line_index INTEGER NOT NULL,
    span_index INTEGER NOT NULL,
    text       TEXT    NOT NULL DEFAULT '',
    style      TEXT    NOT NULL DEFAULT 'default',
    PRIMARY KEY (line_index, span_index)
);
`

// ensureSchema creates all tables if they don't exist and inserts the
// schema_version row when starting fresh.
func ensureSchema(db *sql.DB) error {
	_, err := db.Exec(ddl)
	if err != nil {
		return err
	}

	// Insert version row if missing.
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_version").Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		_, err = db.Exec("INSERT INTO schema_version (version) VALUES (?)", schemaVersion)
	}
	return err
}
