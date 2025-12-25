package models

import (
	"github.com/chempik1234/super-danis-library-golang/pkg/types"
)

// Link is the main entity
type Link struct {
	SourceURL SourceURL
	ShortURL  ShortURL
	CreatedAt types.DateTime
}

// GetUniqueIdentifier - required for caching (genericports.GenericCachePort)
func (l Link) GetUniqueIdentifier() string {
	return l.ShortURL.String()
}

// COOL SOLUTION - types for fields in 1 place

// ShortURL - type for models.Link ShortURL field
type ShortURL = types.AnyText // might be set empty to generate later

// SourceURL - type for models.Link SourceURL field
type SourceURL = types.NotEmptyText
