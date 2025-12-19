package dto

import (
	"fmt"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	types2 "github.com/chempik1234/super-danis-library-golang/pkg/types"
)

// CreateLinkBody is a DTO for create endpoint
type CreateLinkBody struct {
	SourceURL string `json:"source_url"`
	ShortURL  string `json:"short_url,omitempty"`
}

// ToEntity is a method that converts DTO into create-able model (without ID)
func (b CreateLinkBody) ToEntity() (*models.Link, error) {
	var err error

	sourceURL, err := types2.NewNotEmptyText(b.SourceURL)
	if err != nil {
		return nil, fmt.Errorf("source_url musnt't be empty")
	}

	shortURL, err := types2.NewNotEmptyText(b.ShortURL)
	if err != nil {
		return nil, fmt.Errorf("short_url musnt't be empty")
	}

	return &models.Link{
		SourceURL: sourceURL,
		ShortURL:  shortURL,
	}, nil
}
