package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type SQLiteStore struct {
	db *sql.DB
}

var _ Store = (*SQLiteStore)(nil)

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) CreateChop(name, color string, ordered bool) error {
	_, err := s.db.Exec(`
		INSERT INTO chop (name, main_color, ordered)
		VALUES (?, ?, ?)
	`, name, color, ordered)

	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
		return ErrAlreadyExists
	}
	return err
}

func (s *SQLiteStore) AddItem(chop string, item ListItem) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var chopID int64
	err = tx.QueryRow(`SELECT id FROM chop WHERE name = ?`, chop).Scan(&chopID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	var position int
	err = tx.QueryRow(`
		SELECT COALESCE(MAX(position), 0) + 1
		FROM items
		WHERE chop_id = ?
	`, chopID).Scan(&position)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO items (chop_id, text, text_color, description, description_color, position)
		VALUES (?, ?, ?, ?, ?, ?)
	`, chopID, item.Text, item.TextColor, item.Description, item.DescriptionColor, position)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SQLiteStore) Get(name string) (Chop, error) {
	var chop Chop
	err := s.db.QueryRow(`
		SELECT id, name, main_color, ordered
		FROM chop
		WHERE name = ?
	`, name).Scan(&chop.ID, &chop.Name, &chop.MainColor, &chop.Ordered)
	if errors.Is(err, sql.ErrNoRows) {
		return Chop{}, ErrNotFound
	}
	if err != nil {
		return Chop{}, err
	}

	rows, err := s.db.Query(`
		SELECT id, text, text_color, description, description_color, position
		FROM items
		WHERE chop_id = ?
		ORDER BY position
	`, chop.ID)
	if err != nil {
		return Chop{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var it ListItem
		if err := rows.Scan(&it.ID, &it.Text, &it.TextColor, &it.Description, &it.DescriptionColor, &it.Position); err != nil {
			return Chop{}, err
		}
		chop.Items = append(chop.Items, it)
	}
	if err := rows.Err(); err != nil {
		return Chop{}, err
	}
	return chop, nil
}

func (s *SQLiteStore) List(sort SortField) ([]Chop, error) {
	query := `
		SELECT id, name, main_color, ordered
		FROM chop
	`
	switch sort {
	case "", SortByName:
		query += " ORDER BY name"
	case SortByColor:
		query += " ORDER BY main_color"
	case SortByOrdered:
		query += " ORDER BY ordered"
	default:
		return nil, fmt.Errorf("invalid sort %q (valid: name, color, ordered)", sort)
	}
	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chops := []Chop{}
	for rows.Next() {
		var chop Chop
		if err := rows.Scan(&chop.ID, &chop.Name, &chop.MainColor, &chop.Ordered); err != nil {
			return nil, err
		}
		chops = append(chops, chop)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return chops, nil
}

func (s *SQLiteStore) DeleteChop(name string) error {
	res, err := s.db.Exec(`DELETE FROM chop WHERE name = ?`, name)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLiteStore) RmItem(chop string, position int) error {
	res, err := s.db.Exec(`
		DELETE FROM items
		WHERE chop_id = (SELECT id FROM chop WHERE name = ?)
		AND position = ?
	`, chop, position)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func Open(path string) (Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, err
	}
	if err := initSchema(db); err != nil {
		db.Close()
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func initSchema(db *sql.DB) error {

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS chop (
	    id INTEGER PRIMARY KEY,
	    name TEXT NOT NULL UNIQUE,
	    main_color TEXT NOT NULL,
	    ordered INTEGER NOT NULL DEFAULT 0
	);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items (
	    id INTEGER PRIMARY KEY,
		chop_id INTEGER NOT NULL,
	    text TEXT NOT NULL,
	    text_color TEXT NOT NULL DEFAULT 'purple',
		description TEXT NOT NULL DEFAULT '',
		description_color TEXT NOT NULL DEFAULT 'white',
	    position INTEGER NOT NULL,
		UNIQUE (chop_id, position),
		FOREIGN KEY (chop_id) REFERENCES chop(id) ON DELETE CASCADE
	);`)

	return err
}
