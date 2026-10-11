package store

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

var _ Store = (*SQLiteStore)(nil)

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) CreateChop(name, color string, ordered bool) error {
	return errors.New("not implemented")
}

func (s *SQLiteStore) AddItem(chop string, item ListItem) error {
	return errors.New("not implemented")
}

func (s *SQLiteStore) Get(name string) (Chop, error) {
	return Chop{}, errors.New("not implemented")
}

func (s *SQLiteStore) List() ([]Chop, error) {
	return nil, errors.New("not implemented")
}

func (s *SQLiteStore) DeleteChop(name string) error {
	return errors.New("not implemented")
}

func Open(path string) (Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := initSchema(db); err != nil {
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
