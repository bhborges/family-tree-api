package app

import "errors"

var (
	// ErrPersonNotFound occurs when a person is not found.
	ErrPersonNotFound       = errors.New("person not found")
	ErrRelationshipNotFound = errors.New("relationship not found")
	ErrNoRowsInserted       = errors.New("no rows delete")
	ErrNoRowsUpdated        = errors.New("no rows delete")

	// IncestuousOffspring practice not advisable, only for didactic purposes.
	ErrIncestuousOffspring = errors.New("this relationship is not allowed")
)
