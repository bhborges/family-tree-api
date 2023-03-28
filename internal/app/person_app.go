// Package app holds the business logic.
package app

import (
	"context"
	"fmt"

	"github.com/bhborges/family-tree-api/internal/domain"

	"github.com/newrelic/go-agent/v3/newrelic"
)

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
