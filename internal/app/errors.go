package app

import "errors"

var (
	// ErrPersonNotFound occurs when a persons is not found.
	ErrPersonNotFound = errors.New("person not found")
)
