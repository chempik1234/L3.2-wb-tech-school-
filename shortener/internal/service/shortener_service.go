package service

import (
	"context"
	"fmt"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/ports"
	"github.com/chempik1234/super-danis-library-golang/pkg/services"
	"github.com/wb-go/wbf/zlog"
)

// ShortenerService - the service entity that contains business logic related to creating/fetching links
//
// Analytics are also implemented here not to make things complex
type ShortenerService struct {
	cacheService   *services.CachePopularService[string, models.Link]
	storageService ports.ShortenerStorageRepository
}

// NewShortenerService - create new ShortenerService (provide cache service and storage adapter)
func NewShortenerService(
	storage ports.ShortenerStorageRepository,
	cache *services.CachePopularService[string, models.Link]) *ShortenerService {
	return &ShortenerService{
		storageService: storage,
		cacheService:   cache,
	}
}

// CreateLink - create new object in storage
//
// cacheService.minUses < 1 ==> also cache
func (s *ShortenerService) CreateLink(ctx context.Context, model *models.Link) (*models.Link, error) {
	result, err := s.storageService.CreateObject(ctx, model)
	if err != nil {
		return nil, fmt.Errorf("store object: %w", err)
	}

	// cache instantly - no need to wait
	if s.cacheService.MinUsesBeforeCaching() < 1 {
		// cache in background, so Save+Cache and plain Save take equal time
		go func() {
			errCache := s.cacheService.ForceSave(ctx, *model)
			if errCache != nil {
				zlog.Logger.Error().Err(errCache).Msg("error caching after saving")
			}
		}()
	}

	return result, nil
}
