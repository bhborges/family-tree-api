// Package app holds the business logic.
package app

import (
	"context"
	"fmt"

	"github.com/bhborges/family-tree-api/internal/domain"

	"github.com/newrelic/go-agent/v3/newrelic"

	"go.uber.org/zap"
)

// [application name]:[layer].
const _SegmentPrefix = "familytree:application"

// Application definition.
type Application struct {
	repository Repository
	log        *zap.Logger
}

// Repository specifies the signature of a person repository.
type Repository interface {
	GetPersonByID(ctx context.Context, id string) (*domain.Person, error)
	// CreatePerson(ctx context.Context, p domain.Person) (string, error)
	// UpdatePerson(ctx context.Context, p domain.Person) (string, error)
	// DeletePerson(ctx context.Context, p domain.Person) error
}

// NewApplication initializes an instance of a person Application.
func NewApplication(repository Repository, log *zap.Logger) *Application {
	return &Application{repository, log}
}

// GetPersonByID returns a person.
func (a *Application) GetPersonByID(ctx context.Context, id string) (*domain.Person, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "GetPerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	p, err := a.repository.GetPersonByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return p, nil
}
