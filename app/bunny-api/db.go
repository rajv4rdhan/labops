package main

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
CREATE TABLE IF NOT EXISTS bunnies (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	color TEXT NOT NULL DEFAULT 'linen',
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`

	if _, err := s.pool.Exec(ctx, schema); err != nil {
		return err
	}

	var count int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM bunnies`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		_, err := s.pool.Exec(ctx, `INSERT INTO bunnies (name, color) VALUES ('Clover', 'linen'), ('Pip', 'salmon'), ('Mochi', 'umber')`)
		return err
	}
	return nil
}

func (s *Store) ListBunnies(ctx context.Context) ([]Bunny, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, color, created_at FROM bunnies ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bunnies := make([]Bunny, 0)
	for rows.Next() {
		var bunny Bunny
		var createdAt time.Time
		if err := rows.Scan(&bunny.ID, &bunny.Name, &bunny.Color, &createdAt); err != nil {
			return nil, err
		}
		bunny.CreatedAt = createdAt.Format(time.RFC3339)
		bunnies = append(bunnies, bunny)
	}
	return bunnies, rows.Err()
}

func (s *Store) GetBunny(ctx context.Context, id int64) (Bunny, error) {
	var bunny Bunny
	var createdAt time.Time
	err := s.pool.QueryRow(ctx, `SELECT id, name, color, created_at FROM bunnies WHERE id = $1`, id).
		Scan(&bunny.ID, &bunny.Name, &bunny.Color, &createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Bunny{}, err
		}
		return Bunny{}, err
	}
	bunny.CreatedAt = createdAt.Format(time.RFC3339)
	return bunny, nil
}

func (s *Store) CreateBunny(ctx context.Context, req createBunnyRequest) (Bunny, error) {
	color := req.Color
	if color == "" {
		color = "linen"
	}

	var bunny Bunny
	var createdAt time.Time
	err := s.pool.QueryRow(ctx,
		`INSERT INTO bunnies (name, color) VALUES ($1, $2) RETURNING id, name, color, created_at`,
		req.Name, color,
	).Scan(&bunny.ID, &bunny.Name, &bunny.Color, &createdAt)
	if err != nil {
		return Bunny{}, err
	}
	bunny.CreatedAt = createdAt.Format(time.RFC3339)
	return bunny, nil
}
