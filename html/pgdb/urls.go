package pgdb

import (
	"context"

	"github.com/egonelbre/encore-example/go-pure/url"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLs struct {
	db *pgxpool.Pool
}

// URLs returns the urls database.
func (r *Root) URLs() *URLs {
	return &URLs{db: r.db}
}

// Get retrieves url with the specific ID.
func (urls *URLs) Get(ctx context.Context, id url.ID) (*url.URL, error) {
	u := &url.URL{ID: id}
	err := urls.db.QueryRow(ctx, `
		SELECT original_url FROM url
		WHERE id = @id
	`, pgx.NamedArgs{
		"id": u.ID,
	}).Scan(&u.URL)
	return u, err
}

// Insert inserts a URL into the database.
func (urls *URLs) Insert(ctx context.Context, u url.URL) error {
	_, err := urls.db.Exec(ctx, `
		INSERT INTO url (id, original_url)
		VALUES (@id, @url)
	`, pgx.NamedArgs{
		"id":  u.ID,
		"url": u.URL,
	})
	return err
}

// ListAll retrieves all shortened URLs.
func (urls *URLs) ListAll(ctx context.Context) ([]*url.URL, error) {
	rows, err := urls.db.Query(ctx, `
		SELECT id, original_url FROM url
		ORDER BY id
	`)
	result, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*url.URL, error) {
		u := &url.URL{}
		if err := rows.Scan(&u.ID, &u.URL); err != nil {
			return nil, err
		}
		return u, nil
	})
	return result, err
}
