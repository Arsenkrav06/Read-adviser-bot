package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"read-adviser-bot/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?);`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return fmt.Errorf("failed to save page: %w", err)
	}
	return nil
}

func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {

	q := `SELECT url, user_name FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1;`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}
	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?;`
	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return fmt.Errorf("failed to remove page: %w", err)
	}
	return nil
}
func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to check page existence: %w", err)
	}
	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url, user_name TEXT);`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}
