// Service url takes URLs, generates random short IDs, and stores the URLs in a database.
package url

import (
	"context"
	"crypto/rand"
	"encoding/base64"
)

type ID string

type URL struct {
	ID  ID     // short-form URL id
	URL string // complete URL, in long form
}

type ShortenParams struct {
	URL string // the URL to shorten
}

type DB interface {
	Get(ctx context.Context, id ID) (*URL, error)
	Insert(ctx context.Context, u URL) error
	ListAll(ctx context.Context) ([]*URL, error)
}

type Service struct {
	db DB
}

func NewService(db DB) *Service {
	return &Service{db: db}
}

// Shorten shortens a URL.
func (s *Service) Shorten(ctx context.Context, p *ShortenParams) (*URL, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	}

	url := URL{
		ID:  id,
		URL: p.URL,
	}

	err = s.db.Insert(ctx, url)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

// generateID generates a random short ID.
func generateID() (ID, error) {
	var data [6]byte // 6 bytes of entropy
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return ID(base64.RawURLEncoding.EncodeToString(data[:])), nil
}

type ListResponse struct {
	URLs []*URL
}

// List retrieves all shortened URLs.
func (s *Service) List(ctx context.Context) (*ListResponse, error) {
	urls, err := s.db.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return &ListResponse{URLs: urls}, nil
}

// Get retrieves the original URL for the id.
func (s *Service) Get(ctx context.Context, id ID) (*URL, error) {
	return s.db.Get(ctx, id)
}
