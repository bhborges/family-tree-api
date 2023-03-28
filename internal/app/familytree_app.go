package app

import (
	"context"
	"fmt"

	"github.com/bhborges/family-tree-api/internal/domain"

	"github.com/newrelic/go-agent/v3/newrelic"
)

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
