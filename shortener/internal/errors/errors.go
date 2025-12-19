package errors

import "errors"

// ErrLinkNotFound occurs when requested link couldn't be found
//
// Used by both service and repo
var ErrLinkNotFound = errors.New("link not found")

// ErrLinkAlreadyExists occurs when  trying to create with shortURL that already exists
//
// Used by both service and repo
var ErrLinkAlreadyExists = errors.New("shortURL already exists")
