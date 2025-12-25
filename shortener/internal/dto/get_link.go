package dto

import (
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"time"
)

// GetLinkBody is a DTO for getting link information in storage
type GetLinkBody struct {
	SourceURL string `json:"source_url"`
	ShortURL  string `json:"short_url"`
	CreatedAt string `json:"created_at"`
}

// GetLinkBodyToEntity is a method that converts created model to serializable DTO
func GetLinkBodyToEntity(m *models.Link) GetLinkBody {
	return GetLinkBody{
		SourceURL: m.SourceURL.String(),
		ShortURL:  m.ShortURL.String(),
		CreatedAt: m.CreatedAt.Value().Format(time.RFC3339),
	}
}
