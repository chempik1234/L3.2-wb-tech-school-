package analytics

import (
	"context"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
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

// SaveRedirectsBatch - save a bunch of redirects to DB
//
// batching is fast, batching is everything!
func (s StoragePostgresRepo) SaveRedirectsBatch(ctx context.Context, redirect []*models.Redirect) error {
	//TODO implement me
	panic("implement me")
}

// GetAnalytics - get aggregated analytics from inside the DB
//
// Group By is better than a local Golang function
func (s StoragePostgresRepo) GetAnalytics(ctx context.Context, shortLink models.ShortURL) (*models.RedirectDataList, error) {
	//TODO implement me
	panic("implement me")
}
