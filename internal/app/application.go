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
	ListPeople(context.Context) ([]*domain.Person, error)
	GetPersonByID(context.Context, string) (*domain.Person, error)
	CreatePerson(context.Context, domain.Person) (string, error)
	UpdatePerson(context.Context, domain.Person) error
	DeletePerson(context.Context, string) error
	CreateRelationship(context.Context, domain.Relationship) (string, error)
	BuildFamilyTree(context.Context, string) (*domain.FamilyTree, error)
}

// NewApplication initializes an instance of a person Application.
func NewApplication(repository Repository, log *zap.Logger) *Application {
	return &Application{repository, log}
}

// ListPeople return list of people.
func (a *Application) ListPeople(ctx context.Context) ([]*domain.Person, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "ListPeople")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	p, err := a.repository.ListPeople(ctx)
	if err != nil {
		return nil, err
	}

	return p, nil
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
func (a *Application) CreatePerson(ctx context.Context, dp domain.Person) (string, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "CreatePerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	id, err := a.repository.CreatePerson(ctx, dp)
	if err != nil {
		return id, err
	}

	return id, nil
}

// UpdatePerson update a person.
func (a *Application) UpdatePerson(ctx context.Context, dp domain.Person) error {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "UpdatePerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	err := a.repository.UpdatePerson(ctx, dp)

	return err
}

// DeletePerson delete a person.
func (a *Application) DeletePerson(ctx context.Context, id string) error {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "DeletePerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	err := a.repository.DeletePerson(ctx, id)

	return err
}

// CreatePerson create a new person.
func (a *Application) CreateRelationship(ctx context.Context, dr domain.Relationship) (string, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "CreateRelationship")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	id, err := a.repository.CreateRelationship(ctx, dr)
	if err != nil {
		return id, err
	}

	return id, nil
}

// BuildFamilyTree return family tree of person.
func (a *Application) BuildFamilyTree(ctx context.Context, id string) (*domain.FamilyTree, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "BuildFamilyTree")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	t, err := a.repository.BuildFamilyTree(ctx, id)
	if err != nil {
		return nil, err
	}

	return t, nil
}
