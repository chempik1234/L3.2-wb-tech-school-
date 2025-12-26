package ports

import (
	"context"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
)

// ShortenerStorageRepository - port for persistent storage of links
type ShortenerStorageRepository interface {

	// GetObjects - Get all links list from DB
	GetObjects(ctx context.Context) ([]*models.Link, error)

	// GetObjectByID - Get link by shortURL
	//
	// errors.ErrLinkNotFound if not found
	GetObjectByID(ctx context.Context, shortURL models.ShortURL) (*models.Link, error)

	// CreateObject - Create link with given shortURL
	//
	// errors.ErrLinkAlreadyExists if already exists
	//
	// MUTATES object -- sets created_at
	CreateObject(ctx context.Context, fullyReadyObject *models.Link) (*models.Link, error)

	// ObjectExists - check if object with given ID actually exists
	ObjectExists(ctx context.Context, shortURL models.ShortURL) (bool, error)
}

// AnalyticsStorageRepository - port for analytics storage. Save a redirect and get aggregated analytics
//
// Perhaps you'll use the same impl as ShortenerStorageRepository
//
// e.g. everything in Postgres
type AnalyticsStorageRepository interface {
	// SaveRedirectsBatch - save new redirect information
	//
	// They come thousands per second, so do batching please
	SaveRedirectsBatch(ctx context.Context, redirect []*models.Redirect) error

	// GetAnalytics - get aggregated analytics for models.RedirectDataList
	//
	// RESULT IS HALF EMPTY because it can't query LINK model
	//
	// LINK FIELD IS EMPTY QUERY IT YOURSELF with ShortenerStorageRepository
	GetAnalytics(ctx context.Context, shortLink models.ShortURL) (*models.RedirectDataList, error)
}
