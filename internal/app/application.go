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
	CreatePeople(context.Context, []domain.Person) ([]string, error)
	UpdatePerson(context.Context, domain.Person) error
	DeletePerson(context.Context, string) error
	CreateRelationship(context.Context, domain.Relationship) (string, error)
	CreateRelationships(context.Context, []domain.Relationship) ([]string, error)
	UpdateRelationship(context.Context, *domain.Relationship) error
	DeleteRelationship(context.Context, string) error
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

// CreatePeople creates multiple persons.
func (a *Application) CreatePeople(ctx context.Context, people []domain.Person) ([]string, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "CreatePeople")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	personIDs, err := a.repository.CreatePeople(ctx, people)
	if err != nil {
		return personIDs, err
	}

	return personIDs, nil
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

// CreateRelationship create a new relationship.
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

// CreateRelationships creates multiple new relationships.
func (a *Application) CreateRelationships(ctx context.Context, drs []domain.Relationship) ([]string, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "CreateRelationships")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	ids := make([]string, len(drs))
	for i, dr := range drs {
		id, err := a.repository.CreateRelationship(ctx, dr)
		if err != nil {
			return ids, err
		}
		ids[i] = id
	}

	return ids, nil
}

// UpdateRelationship updates an existing relationship.
func (a *Application) UpdateRelationship(ctx context.Context, dr domain.Relationship) error {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "UpdateRelationship")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	err := a.repository.UpdateRelationship(ctx, &dr)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRelationship deletes a relationship.
func (a *Application) DeleteRelationship(ctx context.Context, id string) error {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "DeleteRelationship")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	err := a.repository.DeleteRelationship(ctx, id)
	if err != nil {
		return err
	}

	return nil
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
