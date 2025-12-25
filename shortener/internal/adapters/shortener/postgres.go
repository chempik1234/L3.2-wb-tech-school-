package shortener

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

// GetObjects - Get all links list from DB
func (s *StoragePostgresRepo) GetObjects(ctx context.Context) ([]*models.Link, error) {
	//TODO implement me
	panic("implement me")
}

// GetObjectByID - Get link by shortURL
//
// errors.ErrLinkNotFound if not found
func (s *StoragePostgresRepo) GetObjectByID(ctx context.Context, shortURL string) (*models.Link, error) {
	//TODO implement me
	panic("implement me")
}

// CreateObject - Create link with given shortURL
//
// errors.ErrLinkAlreadyExists if already exists
func (s *StoragePostgresRepo) CreateObject(ctx context.Context, fullyReadyObject *models.Link) (*models.Link, error) {
	//TODO implement me
	panic("implement me")
}

// UpdateObject - Update link with given shortURL in it
//
// errors.ErrLinkNotFound if not found
func (s *StoragePostgresRepo) UpdateObject(ctx context.Context, fullyReadyObject *models.Link) (*models.Link, error) {
	//TODO implement me
	panic("implement me")
}

// DeleteObject - delete link with given shortURL
//
// errors.ErrLinkNotFound if not found
func (s *StoragePostgresRepo) DeleteObject(ctx context.Context, shortURL string) error {
	//TODO implement me
	panic("implement me")
}

// ObjectExists - delete link with given shortURL
//
// errors.ErrLinkNotFound if not found
func (s *StoragePostgresRepo) ObjectExists(ctx context.Context, shortURL string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
