package app

import (
	"context"
	"fmt"

	"github.com/bhborges/family-tree-api/internal/domain"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// ListRelationships list all relationships
func (a *Application) ListRelationships(ctx context.Context) ([]*domain.Relationship, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "ListRelationships")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	p, err := a.repository.ListRelationships(ctx)
	if err != nil {
		return nil, err
	}

	return p, nil
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
