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
	GetPersonByID(context.Context, string) (*domain.Person, error)
	CreatePerson(context.Context, domain.Person) (string, error)
	UpdatePerson(context.Context, domain.Person) (error)
	DeletePerson(context.Context, string) error
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

// CreatePerson create a new person.
func (a *Application) CreatePerson(ctx context.Context, p domain.Person) (string, error)  {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "CreatePerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	id, err := a.repository.CreatePerson(ctx, p)
	if err != nil {
		return id, err
	}

	return id, nil
}

// UpdatePerson update a person.
func (a *Application) UpdatePerson(ctx context.Context, p domain.Person) (error)  {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "UpdatePerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	err := a.repository.UpdatePerson(ctx, p)

	return err
}

// DeletePerson delete a person.
func (a *Application) DeletePerson(ctx context.Context, id string) (error)  {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "DeletePerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	err := a.repository.DeletePerson(ctx, id)

	return err
}