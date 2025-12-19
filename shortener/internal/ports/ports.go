package ports

import (
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"github.com/chempik1234/super-danis-library-golang/pkg/genericports"
)

// ShortenerStorageRepository - port for persistent storage of links
type ShortenerStorageRepository = genericports.GenericStoragePort[string, models.Link]
