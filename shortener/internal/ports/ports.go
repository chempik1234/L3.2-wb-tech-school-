package ports

import (
	"context"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"github.com/chempik1234/super-danis-library-golang/pkg/genericports"
)

// ShortenerStorageRepository - port for persistent storage of links
type ShortenerStorageRepository = objectStorageEnhanced[string, models.Link]

type objectStorageEnhanced[I comparable, V genericports.ObjectWithIdentifier[I]] interface {
	genericports.GenericStoragePort[I, V]
	// ObjectExists - check if object with given ID actually exists
	ObjectExists(ctx context.Context, id I) (bool, error)
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
	GetAnalytics(ctx context.Context, shortLink models.ShortURL) (*models.RedirectDataList, error)
}
