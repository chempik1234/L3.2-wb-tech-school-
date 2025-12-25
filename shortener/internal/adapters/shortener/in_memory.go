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
func (s *StorageInMemoryRepo) GetObjects(ctx context.Context) ([]*models.Link, error) {
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
func (s *StorageInMemoryRepo) GetObjectByID(ctx context.Context, id string) (*models.Link, error) {
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
func (s *StorageInMemoryRepo) CreateObject(ctx context.Context, fullyReadyObject *models.Link) (*models.Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[fullyReadyObject.GetUniqueIdentifier()]; exists {
		return nil, errors.ErrLinkAlreadyExists
	}

	s.data[fullyReadyObject.GetUniqueIdentifier()] = fullyReadyObject
	return fullyReadyObject, nil
}

// UpdateObject updates an existing Link object in storage
//
// uses shortURL given in "fullyReadyObject" to identify the record being edited
func (s *StorageInMemoryRepo) UpdateObject(ctx context.Context, fullyReadyObject *models.Link) (*models.Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[fullyReadyObject.GetUniqueIdentifier()]; !exists {
		return nil, errors.ErrLinkNotFound
	}

	s.data[fullyReadyObject.GetUniqueIdentifier()] = fullyReadyObject
	return fullyReadyObject, nil
}

// DeleteObject removes a Link object from storage by its ID
//
// err if not exists
func (s *StorageInMemoryRepo) DeleteObject(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[id]; !exists {
		return errors.ErrLinkNotFound
	}

	delete(s.data, id)
	return nil
}
