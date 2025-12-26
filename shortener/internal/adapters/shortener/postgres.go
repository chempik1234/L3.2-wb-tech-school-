package shortener

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/adapters"
	errors2 "github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/errors"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"github.com/chempik1234/super-danis-library-golang/pkg/types"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
	"time"
)

// StoragePostgresRepo - adapter for ports.StoragePostgresRepo
//
// PostgresSQL
type StoragePostgresRepo struct {
	db       *dbpg.DB
	strategy retry.Strategy
}

// NewStoragePostgresRepo creates a new StoragePostgresRepo
func NewStoragePostgresRepo(db *dbpg.DB, retryStrategy retry.Strategy) *StoragePostgresRepo {
	return &StoragePostgresRepo{db: db, strategy: retryStrategy}
}

// GetObjects - Get all links list from DB
func (s *StoragePostgresRepo) GetObjects(ctx context.Context) ([]*models.Link, error) {
	query := `SELECT short_url, source_url, created_at FROM links`
	rows, err := s.db.QueryWithRetry(ctx, s.strategy, query)
	if err != nil {
		return nil, fmt.Errorf("error selecting all rows: %w", err)
	}

	defer adapters.ClosePostgresRows(rows)
	links := make([]*models.Link, 0)
	for rows.Next() {
		link := &models.Link{}
		err = rows.Scan(&link.ShortURL, &link.SourceURL, &link.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		links = append(links, link)
	}

	return links, nil
}

// GetObjectByID - Get link by shortURL
//
// errors.ErrLinkNotFound if not found
func (s *StoragePostgresRepo) GetObjectByID(ctx context.Context, shortURL models.ShortURL) (*models.Link, error) {
	query := `SELECT short_url, source_url, created_at FROM links WHERE short_url = $1`
	row, err := s.db.QueryRowWithRetry(ctx, s.strategy, query, shortURL.String())
	if err != nil {
		return nil, fmt.Errorf("error selecting row: %w", err)
	}

	createdAt := time.Time{}

	link := &models.Link{}
	err = row.Scan(&link.ShortURL, &link.SourceURL, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors2.ErrLinkNotFound
		}
		return nil, fmt.Errorf("error scanning row: %w", err)
	}

	link.CreatedAt = types.NewDateTime(createdAt)

	return link, nil
}

// CreateObject - Create link with given shortURL
//
// errors.ErrLinkAlreadyExists if already exists
//
// MUTATES object -- sets created_at
func (s *StoragePostgresRepo) CreateObject(ctx context.Context, fullyReadyObject *models.Link) (*models.Link, error) {
	query := `INSERT INTO links (source_url, short_url)
				VALUES ($1, $2)
				ON CONFLICT (short_url) DO NOTHING
				RETURNING created_at` // let's NOT create a separate schema for our tables
	row, err := s.db.QueryRowWithRetry(ctx, s.strategy, query, fullyReadyObject.SourceURL, fullyReadyObject.ShortURL)
	if err != nil {
		return nil, fmt.Errorf("error querying postgres after retries: %w", err)
	}

	createdAt := time.Time{}

	err = row.Scan(&createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors2.ErrLinkAlreadyExists
		}
		return nil, fmt.Errorf("unknown error scanning row: %w", err)
	}

	fullyReadyObject.CreatedAt = types.NewDateTime(createdAt)

	return fullyReadyObject, nil
}

// ObjectExists - delete link with given shortURL
//
// errors.ErrLinkNotFound if not found
func (s *StoragePostgresRepo) ObjectExists(ctx context.Context, shortURL models.ShortURL) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM links WHERE short_url = $1)`
	row, err := s.db.QueryRowWithRetry(ctx, s.strategy, query, shortURL.String())
	if err != nil {
		return false, fmt.Errorf("error checking if row exists: %w", err)
	}

	exists := false
	err = row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error scanning row: %w", err)
	}

	return exists, nil
}
