package service

import (
	"context"
	"fmt"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/ports"
	"github.com/chempik1234/super-danis-library-golang/pkg/services"
	"github.com/chempik1234/super-danis-library-golang/pkg/types"
	"github.com/wb-go/wbf/zlog"
	"math/rand"
	"time"
)

const batchingChannelSize = 1000

// ShortenerService - the service entity that contains business logic related to creating/fetching links
//
// # Analytics are also implemented here not to make things complex
type ShortenerService struct {
	cacheService               *services.CachePopularService[string, models.Link]
	shortenerStorageRepository ports.ShortenerStorageRepository
	analyticsStorageRepository ports.AnalyticsStorageRepository

	maxLinkLen     int
	batchingPeriod time.Duration

	// init chan only when running saving in background!
	redirectsForBatching chan *models.Redirect
}

// NewShortenerService - create new ShortenerService (provide cache service and storage adapter)
func NewShortenerService(
	shortenerStorage ports.ShortenerStorageRepository,
	analyticsStorage ports.AnalyticsStorageRepository,
	cache *services.CachePopularService[string, models.Link],
	maxLinkLen int,
	batchingPeriod time.Duration,
) *ShortenerService {
	return &ShortenerService{
		shortenerStorageRepository: shortenerStorage,
		analyticsStorageRepository: analyticsStorage,
		cacheService:               cache,
		maxLinkLen:                 maxLinkLen,
		redirectsForBatching:       nil, // init channel only in Run...
		batchingPeriod:             batchingPeriod,
	}
}

// CreateLink - create new object in storage
//
// cacheService.minUses < 1 ==> also cache
func (s *ShortenerService) CreateLink(ctx context.Context, model *models.Link) (*models.Link, error) {
	if len(model.ShortURL.String()) == 0 {
		link := s.generateURL(ctx, s.maxLinkLen)
		model.ShortURL = link
	}

	result, err := s.shortenerStorageRepository.CreateObject(ctx, model)
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

// GetLink - get link by id (shortLink)
//
// Use to check if link exists before redirect
func (s *ShortenerService) GetLink(ctx context.Context, linkString models.ShortURL) (*models.Link, error) {
	var link *models.Link
	var err error
	// step 1. try to get from cache
	if link, err = s.cacheService.Get(ctx, linkString.String()); err != nil {
		// step 2. try to get from storage
		if link, err = s.shortenerStorageRepository.GetObjectByID(ctx, linkString.String()); err != nil {
			return nil, fmt.Errorf("storage error: %w", err)
		}

		go func() {
			errCache := s.cacheService.UpdatePopularity(ctx, *link, 1)
			if errCache != nil {
				zlog.Logger.Error().Err(errCache).Stringer("short_url", link.ShortURL).Msg("error caching after saving popularity for shortURL")
			}
		}()
	}

	return link, err
}

// SaveRedirect - creates record in analytics table
//
// WORKS ONLY after ShortenerService.RunBatchSavingInBackground has started!
func (s *ShortenerService) SaveRedirect(ctx context.Context, shortLink models.ShortURL, userAgent types.AnyText, clickAt types.DateTime) error {
	if s.redirectsForBatching != nil {
		s.redirectsForBatching <- &models.Redirect{
			ClickAt:   clickAt,
			UserAgent: userAgent,
			ShortURL:  shortLink,
		}
	}
	return nil
}

// GetAnalytics - return aggregated models.RedirectDataList analytics
func (s *ShortenerService) GetAnalytics(ctx context.Context, shortURL models.ShortURL) (*models.RedirectDataList, error) {
	data, err := s.analyticsStorageRepository.GetAnalytics(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("analytics error: %w", err)
	}
	return data, nil
}

func (s *ShortenerService) generateURL(ctx context.Context, length int) types.AnyText {
	// loop before we get unique link
	var err error
	var alreadyExists bool

	link := generateRandomString(length)

	for {
		if alreadyExists, err = s.LinkExists(ctx, link); !alreadyExists && err == nil {
			break
		}
		if err != nil {
			zlog.Logger.Error().Err(err).Stringer("link", link).Msg("error checking if link exists")
		}
		link = generateRandomString(length)
	}

	return link
}

// LinkExists - check if link actually exists in cache or storage
func (s *ShortenerService) LinkExists(ctx context.Context, link types.AnyText) (bool, error) {
	return s.shortenerStorageRepository.ObjectExists(ctx, link.String())
}

// RunBatchSavingInBackground - run saving redirects in batches
//
// # Every T period it reads channel and saves tons of redirects
//
// 95% of T - read from chan
//
//	// after that
//	go save()
//
// Stops on ctx.Done()
func (s *ShortenerService) RunBatchSavingInBackground(ctx context.Context) {
	if s.redirectsForBatching == nil {
		// on start, init chan
		s.redirectsForBatching = make(chan *models.Redirect, batchingChannelSize)
	}

	tickTimer := time.NewTicker(s.batchingPeriod)
	defer tickTimer.Stop()

	sharedBatch := make([]*models.Redirect, 0, batchingChannelSize)

l:
	for {
		select {
		// every T period we do read channel
		case <-tickTimer.C:
			// step 1. reset already allocated slice
			sharedBatch = sharedBatch[:0] // I hope we don't lose anything!

			// step 2. read from chan into slice, at max we read for T/2
			ctxReadFromChan, cancel := context.WithTimeout(ctx, s.batchingPeriod*95/10)

		l2:
			for {
				select {
				case obj := <-s.redirectsForBatching:
					sharedBatch = append(sharedBatch, obj)
				case <-ctxReadFromChan.Done():
					break l2
				}
			}

			cancel()

			// step 3.1. if nothing to save, don't call save
			if len(sharedBatch) == 0 {
				break
			}

			// step 4. allocate batch to save before going to new one
			batchToSave := make([]*models.Redirect, len(sharedBatch))
			copy(batchToSave, sharedBatch)

			// step 5. send to "try to save"
			go s.saveBatch(ctx, batchToSave)

		case <-ctx.Done():
			// on stop, close chan
			redirectsChan := s.redirectsForBatching
			s.redirectsForBatching = nil
			close(redirectsChan)

			break l
		}
	}
}

func (s *ShortenerService) saveBatch(ctx context.Context, save []*models.Redirect) {
	err := s.analyticsStorageRepository.SaveRedirectsBatch(ctx, save)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("error saving redirects batch")
	}
}

func generateRandomString(length int) types.AnyText {
	resultStringBytes := make([]byte, length)

	// random letters a-z
	for i := 0; i < length; i++ {
		resultStringBytes[i] = byte(rand.Intn(25)) + 97
	}

	return types.NewAnyText(string(resultStringBytes))
}
