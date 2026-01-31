// Service url takes URLs, generates random short IDs, and stores the URLs in a database.
package url

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"encore.dev/storage/sqldb"

	"github.com/jackc/pgx/v5"
)

type URL struct {
	ID  string // short-form URL id
	URL string // complete URL, in long form
}

type ShortenParams struct {
	URL string // the URL to shorten
}

// Shorten shortens a URL.
//
//encore:api public method=POST path=/api/url
func Shorten(ctx context.Context, p *ShortenParams) (*URL, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	} else if err := insert(ctx, id, p.URL); err != nil {
		return nil, err
	}
	return &URL{ID: id, URL: p.URL}, nil
}

// generateID generates a random short ID.
func generateID() (string, error) {
	var data [6]byte // 6 bytes of entropy
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

// insert inserts a URL into the database.
func insert(ctx context.Context, id, url string) error {
	_, err := db.Exec(ctx, `
		INSERT INTO url (id, original_url)
		VALUES (@id, @url)
	`, pgx.NamedArgs{
		"id":  id,
		"url": url,
	})
	return err
}

type ListResponse struct {
	URLs []*URL
}

// List retrieves all shortened URLs.
//
//encore:api public method=GET path=/api/url
func List(ctx context.Context) (*ListResponse, error) {
	rows, err := db.Query(ctx, `
		SELECT id, original_url FROM url
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []*URL
	for rows.Next() {
		u := &URL{}
		if err := rows.Scan(&u.ID, &u.URL); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return &ListResponse{URLs: urls}, nil
}

// Get retrieves the original URL for the id.
//
//encore:api public method=GET path=/api/url/:id
func Get(ctx context.Context, id string) (*URL, error) {
	u := &URL{ID: id}
	err := db.QueryRow(ctx, `
		SELECT original_url FROM url
		WHERE id = $1
	`, id).Scan(&u.URL)
	return u, err
}

// Below we define a database named 'url', using the database
// migrations  in the "./migrations" folder.
// Encore provisions, migrates, and connects to the database.
// Learn more: https://encore.dev/docs/go/primitives/databases

// 'url' database is used to store the URLs that are being shortened.
var db = sqldb.NewDatabase("url", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
