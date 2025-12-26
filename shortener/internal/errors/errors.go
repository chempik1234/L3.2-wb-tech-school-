package errors

import (
	"errors"
	"fmt"
)

// ErrLinkNotFound occurs when requested link couldn't be found
//
// Used by both service and repo
var ErrLinkNotFound = errors.New("link not found")

// ErrLinkAlreadyExists occurs when  trying to create with shortURL that already exists
//
// Used by both service and repo
var ErrLinkAlreadyExists = errors.New("shortURL already exists")

// ErrValidation - validation error use with NewValidationError
var ErrValidation = errors.New("validation error")

// NewValidationError - create ErrValidation with custom description
//
// Use this -> auto return http.StatusBadRequest
func NewValidationError(err error) error {
	return fmt.Errorf("%w: %w", ErrValidation, err)
}
