package shortener

import (
	"context"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/errors"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"sync"
)

// StorageInMemoryRepo - impl ports.ShortenerStorageRepository (Map, in-memory)
type StorageInMemoryRepo struct {
	mu   *sync.RWMutex
	data map[string]*models.Link
}

// NewStorageInMemoryRepo - creates new instance of *NewStorageInMemoryRepo.
func NewStorageInMemoryRepo() *StorageInMemoryRepo {
	return &StorageInMemoryRepo{
		data: make(map[string]*models.Link),
		mu:   new(sync.RWMutex),
	}
}

// GetObjects retrieves all Link objects from storage
func (s *StorageInMemoryRepo) GetObjects(_ context.Context) ([]*models.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	links := make([]*models.Link, 0, len(s.data))
	for _, link := range s.data {
		links = append(links, link)
	}
	return links, nil
}

// GetObjectByID retrieves a Link object by its ID
//
// error on not exists
func (s *StorageInMemoryRepo) GetObjectByID(_ context.Context, id string) (*models.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	link, exists := s.data[id]
	if !exists {
		return nil, errors.ErrLinkNotFound
	}
	return link, nil
}

// CreateObject adds a new Link object to storage
//
// Given shortURL is used, so pre-generate it!
//
// Error on conflict
func (s *StorageInMemoryRepo) CreateObject(_ context.Context, fullyReadyObject *models.Link) (*models.Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[fullyReadyObject.GetUniqueIdentifier()]; exists {
		return nil, errors.ErrLinkAlreadyExists
	}

	s.data[fullyReadyObject.GetUniqueIdentifier()] = fullyReadyObject
	return fullyReadyObject, nil
}

func (s *StorageInMemoryRepo) ObjectExists(_ context.Context, shortURL models.ShortURL) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[shortURL.String()]
	return ok, nil
}
