package transport

import (
	"context"
	"fmt"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/dto"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ShortenerHandler is the HTTP routes handler, used in AssembleRouter
//
// Validates request and passes it to service layer
type ShortenerHandler struct {
	crudService *service.ShortenerService
}

// NewShortenerHandler creates a new ShortenerHandler with given service
func NewShortenerHandler(crudService *service.ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{crudService: crudService}
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

	result, err := h.crudService.CreateLink(context.Background(), createModel)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusConflict,
			gin.H{"error": fmt.Sprintf("couldn't perform operation: %s", err.Error())},
		)
		return
	}

	// CreatedAt is now assigned!
	c.JSON(http.StatusCreated, dto.GetLinkBodyToEntity(result))
}
