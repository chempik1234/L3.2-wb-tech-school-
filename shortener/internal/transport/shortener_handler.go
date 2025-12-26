package transport

import (
	"context"
	"errors"
	"fmt"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/dto"
	errors2 "github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/errors"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/service"
	"github.com/chempik1234/super-danis-library-golang/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/zlog"
	"net/http"
	"time"
)

const (
	shortLinkParam = "short_url"
)

// ShortenerHandler is the HTTP routes handler, used in AssembleRouter
//
// Validates request and passes it to service layer
type ShortenerHandler struct {
	shortenerService *service.ShortenerService
}

// NewShortenerHandler creates a new ShortenerHandler with given service
func NewShortenerHandler(crudService *service.ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{shortenerService: crudService}
}

// CreateLink POST /shorten
func (h *ShortenerHandler) CreateLink(c *gin.Context) {
	var body dto.CreateLinkBody

	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid body (parsing): %s", err.Error())},
		)
		return
	}

	var createModel *models.Link
	createModel, err = body.ToEntity()
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid body (validating): %s", err.Error())},
		)
		return
	}

	result, err := h.shortenerService.CreateLink(context.Background(), createModel)
	if err != nil {
		zlog.Logger.Error().Err(err).Any("body", body).Msg("couldn't create link")
		c.AbortWithStatusJSON(
			h.statusForError(err),
			gin.H{"error": fmt.Sprintf("couldn't perform operation: %s", err.Error())},
		)
		return
	}

	// CreatedAt is now assigned!
	c.JSON(http.StatusCreated, dto.GetLinkBodyToEntity(result))
}

// RedirectLink GET /s/:short_url
func (h *ShortenerHandler) RedirectLink(c *gin.Context) {
	shortLink, link, err := h.getShortLinkAndLink(c)
	if err != nil || link == nil {
		c.AbortWithStatusJSON(
			h.statusForError(err),
			gin.H{"error": err},
		)
		return
	}

	userAgent := types.NewAnyText(c.GetHeader("User-Agent"))

	zlog.Logger.Info().Stringer("user_agent", userAgent).Msg("new redirect")

	go func() {
		saveErr := h.shortenerService.SaveRedirect(
			context.Background(),
			// convert to ShortURL because in this case, we must use types.NotEmptyText to validate
			// and not types.AnyText which ShortURL IS under the hood
			models.ShortURL(shortLink),
			userAgent,
			types.NewDateTime(time.Now()),
		)
		if saveErr != nil {
			zlog.Logger.Error().Err(saveErr).Msg("error saving link")
		}
	}()

	c.Redirect(http.StatusPermanentRedirect, link.SourceURL.String())
}

// AnalyticsLink GET /analytics/:short_url
func (h *ShortenerHandler) AnalyticsLink(c *gin.Context) {
	shortLink, link, err := h.getShortLinkAndLink(c)
	if err != nil || link == nil {
		c.AbortWithStatusJSON(
			h.statusForError(err),
			gin.H{"error": err},
		)
	}

	analyticsData, err := h.shortenerService.GetAnalytics(context.Background(), link)
	if err != nil {
		zlog.Logger.Error().Err(err).Stringer(shortLinkParam, shortLink).Msg("couldn't create link")
		c.AbortWithStatusJSON(
			h.statusForError(err),
			gin.H{"error": fmt.Sprintf("couldn't perform operation: %s", err.Error())},
		)
	}

	c.JSON(http.StatusOK, dto.AnalyticsBodyFromDataList(analyticsData))
}

func (h *ShortenerHandler) getShortLinkAndLink(c *gin.Context) (types.NotEmptyText, *models.Link, error) {
	// models.ShortURL is actually types2.NotEmptyText
	shortLink, err := types.NewNotEmptyText(c.Param(shortLinkParam))
	if err != nil {
		return "", nil, errors2.NewValidationError(errors.New("link mustn't be empty"))
	}

	var link *models.Link
	link, err = h.shortenerService.GetLink(context.Background(), models.ShortURL(shortLink))
	if err != nil {
		return shortLink, nil, fmt.Errorf("error getting link for '%s': %w", shortLink, err)
	}

	if link == nil {
		return shortLink, nil, fmt.Errorf("link not found: %s", shortLink)
	}

	return shortLink, link, nil
}

func (h *ShortenerHandler) statusForError(err error) int {
	if errors.Is(err, errors2.ErrLinkNotFound) {
		return http.StatusNotFound
	} else if errors.Is(err, errors2.ErrLinkAlreadyExists) {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}
