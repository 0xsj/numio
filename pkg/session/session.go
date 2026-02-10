// Package session provides SQLite-backed persistence for numio sessions.
// It stores buffer text, cursor position, variables, and result spans so the
// desktop app can restore exactly where the user left off.
package session

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // pure-Go SQLite driver

	"github.com/0xsj/numio/internal/rpc"
	"github.com/0xsj/numio/pkg/types"
)

// SessionData is the transfer type between the editor and the session store.
// The editor produces/consumes this struct; the session package reads/writes
// it to SQLite. No coupling to editor internals.
type SessionData struct {
	Lines     []string
	CursorRow int
	CursorCol int
	Variables map[string]types.Value
	Results   map[int][]rpc.Span
}

// Store manages a single SQLite database file for session persistence.
type Store struct {
	db   *sql.DB
	path string
}

// DefaultPath returns ~/.numio/numio.db.
func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".numio", "numio.db")
}

// Open opens (or creates) the session database at the given path.
// If schema creation fails it deletes the file and retries once.
func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	s := &Store{path: path}
	if err := s.open(); err != nil {
		// Corruption recovery: remove and retry once.
		os.Remove(path)
		if err2 := s.open(); err2 != nil {
			return nil, err2
		}
	}
	return s, nil
}

func (s *Store) open() error {
	db, err := sql.Open("sqlite", s.path)
	if err != nil {
		return err
	}
	// WAL mode for better concurrency.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return err
	}
	if err := ensureSchema(db); err != nil {
		db.Close()
		return err
	}
	s.db = db
	return nil
}

// Close closes the database.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Save persists session data. It does a full replace inside a single
// transaction (DELETE + INSERT). Calculator buffers are small so this is fine.
func (s *Store) Save(data *SessionData) error {
	if s.db == nil || data == nil {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// ── buffer ──────────────────────────────────────────────────
	if _, err := tx.Exec("DELETE FROM buffer"); err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO buffer (row_index, text) VALUES (?, ?)")
	if err != nil {
		return err
	}
	for i, line := range data.Lines {
		if _, err := stmt.Exec(i, line); err != nil {
			stmt.Close()
			return err
		}
	}
	stmt.Close()

	// ── cursor ──────────────────────────────────────────────────
	if _, err := tx.Exec("DELETE FROM cursor"); err != nil {
		return err
	}
	if _, err := tx.Exec(
		"INSERT INTO cursor (id, row, col) VALUES (1, ?, ?)",
		data.CursorRow, data.CursorCol,
	); err != nil {
		return err
	}

	// ── variables ───────────────────────────────────────────────
	if _, err := tx.Exec("DELETE FROM variables"); err != nil {
		return err
	}
	vstmt, err := tx.Prepare(
		"INSERT INTO variables (name, kind, num, str, type_code, period) VALUES (?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	for name, val := range data.Variables {
		kind, num, str, typeCode, period := serializeValue(val)
		if _, err := vstmt.Exec(name, kind, num, str, typeCode, period); err != nil {
			vstmt.Close()
			return err
		}
	}
	vstmt.Close()

	// ── results ─────────────────────────────────────────────────
	if _, err := tx.Exec("DELETE FROM results"); err != nil {
		return err
	}
	rstmt, err := tx.Prepare(
		"INSERT INTO results (line_index, span_index, text, style) VALUES (?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	for lineIdx, spans := range data.Results {
		for spanIdx, span := range spans {
			if _, err := rstmt.Exec(lineIdx, spanIdx, span.Text, string(span.Style)); err != nil {
				rstmt.Close()
				return err
			}
		}
	}
	rstmt.Close()

	return tx.Commit()
}

// Load reads the persisted session. Returns nil if the database is empty or
// if any error occurs (graceful degradation: start fresh).
func (s *Store) Load() *SessionData {
	if s.db == nil {
		return nil
	}

	data := &SessionData{
		Variables: make(map[string]types.Value),
		Results:   make(map[int][]rpc.Span),
	}

	// ── buffer ──────────────────────────────────────────────────
	rows, err := s.db.Query("SELECT row_index, text FROM buffer ORDER BY row_index")
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var idx int
		var text string
		if err := rows.Scan(&idx, &text); err != nil {
			return nil
		}
		// Grow slice to fit.
		for len(data.Lines) <= idx {
			data.Lines = append(data.Lines, "")
		}
		data.Lines[idx] = text
	}
	if rows.Err() != nil {
		return nil
	}

	// Nothing saved yet.
	if len(data.Lines) == 0 {
		return nil
	}

	// ── cursor ──────────────────────────────────────────────────
	err = s.db.QueryRow("SELECT row, col FROM cursor WHERE id = 1").
		Scan(&data.CursorRow, &data.CursorCol)
	if err != nil {
		// Missing cursor row is fine, defaults are 0,0.
		data.CursorRow = 0
		data.CursorCol = 0
	}

	// ── variables ───────────────────────────────────────────────
	vrows, err := s.db.Query("SELECT name, kind, num, str, type_code, period FROM variables")
	if err == nil {
		defer vrows.Close()
		for vrows.Next() {
			var name, str, typeCode string
			var kind, period int
			var num float64
			if err := vrows.Scan(&name, &kind, &num, &str, &typeCode, &period); err != nil {
				continue
			}
			data.Variables[name] = deserializeValue(kind, num, str, typeCode, period)
		}
	}

	// ── results ─────────────────────────────────────────────────
	rrows, err := s.db.Query("SELECT line_index, span_index, text, style FROM results ORDER BY line_index, span_index")
	if err == nil {
		defer rrows.Close()
		for rrows.Next() {
			var lineIdx, spanIdx int
			var text, style string
			if err := rrows.Scan(&lineIdx, &spanIdx, &text, &style); err != nil {
				continue
			}
			data.Results[lineIdx] = append(data.Results[lineIdx], rpc.Span{
				Text:  text,
				Style: rpc.SpanStyle(style),
			})
		}
	}

	return data
}
