package models

import (
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/pkg/types"
	types2 "github.com/chempik1234/super-danis-library-golang/pkg/types"
)

// Link is the main entity
type Link struct {
	SourceURL types2.NotEmptyText
	ShortURL  types2.NotEmptyText
	CreatedAt types.DateTime
}

// GetUniqueIdentifier - required for caching (genericports.GenericCachePort)
func (l Link) GetUniqueIdentifier() string {
	return l.ShortURL.String()
}
