package main

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("note not found")

type Note struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(ctx context.Context, dsn string) (*Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	store := &Store{pool: pool}
	if err := store.migrate(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *Store) migrate(ctx context.Context) error {
	const schema = `
CREATE TABLE IF NOT EXISTS notes (
	id BIGSERIAL PRIMARY KEY,
	title TEXT NOT NULL DEFAULT '',
	body TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`

	if _, err := s.pool.Exec(ctx, schema); err != nil {
		return err
	}

	var count int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM notes`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		_, err := s.pool.Exec(ctx,
			`INSERT INTO notes (title, body) VALUES ($1, $2)`,
			"Welcome to Squirrel 🐿️",
			"Stash your thoughts here. Create a note, write something, and it is saved to Postgres.\n\nThere is no login — every note is shared with everyone.",
		)
		return err
	}
	return nil
}

func (s *Store) ListNotes(ctx context.Context) ([]Note, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, title, body, created_at, updated_at FROM notes ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]Note, 0)
	for rows.Next() {
		note, err := scanNote(rows)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, rows.Err()
}

func (s *Store) GetNote(ctx context.Context, id int64) (Note, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, title, body, created_at, updated_at FROM notes WHERE id = $1`, id)
	note, err := scanNote(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Note{}, ErrNotFound
	}
	return note, err
}

func (s *Store) CreateNote(ctx context.Context, title, body string) (Note, error) {
	row := s.pool.QueryRow(ctx,
		`INSERT INTO notes (title, body) VALUES ($1, $2)
		 RETURNING id, title, body, created_at, updated_at`,
		title, body)
	return scanNote(row)
}

func (s *Store) UpdateNote(ctx context.Context, id int64, title, body string) (Note, error) {
	row := s.pool.QueryRow(ctx,
		`UPDATE notes SET title = $2, body = $3, updated_at = now()
		 WHERE id = $1
		 RETURNING id, title, body, created_at, updated_at`,
		id, title, body)
	note, err := scanNote(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Note{}, ErrNotFound
	}
	return note, err
}

func (s *Store) DeleteNote(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM notes WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanNote(row scanner) (Note, error) {
	var note Note
	var createdAt, updatedAt time.Time
	if err := row.Scan(&note.ID, &note.Title, &note.Body, &createdAt, &updatedAt); err != nil {
		return Note{}, err
	}
	note.CreatedAt = createdAt.Format(time.RFC3339)
	note.UpdatedAt = updatedAt.Format(time.RFC3339)
	return note, nil
}
