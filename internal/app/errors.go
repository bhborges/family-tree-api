package app

import "errors"

var (
	// ErrPersonNotFound occurs when a person is not found.
	ErrPersonNotFound = errors.New("person not found")
)
