// Package adapters implements all the necessary
// logic to talk to any external system.
package adapter

import (
	"context"
	"fmt"

	"github.com/bhborges/family-tree-api/internal/domain"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// CreateRelationship create a new relationship.
func (pr *PostgresRepository) CreateRelationship(ctx context.Context, dr domain.Relationship) (string, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "CreateRelationShip")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	r := domain.Relationship{
		ParentID: dr.ParentID,
		ChildID:  dr.ChildID,
	}

	tx := pr.db.Create(&r)
	if tx.Error != nil {
		return "", tx.Error
	}

	return r.ID, nil
}

func (pr *PostgresRepository) CreateRelationships(ctx context.Context, relationships []domain.Relationship) ([]string, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "CreateRelationships")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	tx := pr.db.Begin()
	ids := make([]string, len(relationships))
	for i, r := range relationships {
		if err := tx.Create(&r).Error; err != nil {
			tx.Rollback()
			return ids, err
		}
		ids[i] = r.ID
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return ids, err
	}

	return ids, nil
}

func (pr *PostgresRepository) UpdateRelationship(ctx context.Context, dr *domain.Relationship) error {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "UpdateRelationShip")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	tx := pr.db.Model(&domain.Relationship{}).
		Where("id = ?", dr.ID).
		Updates(map[string]interface{}{
			"parent_id": dr.ParentID,
			"child_id":  dr.ChildID,
		})

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

func (pr *PostgresRepository) DeleteRelationship(ctx context.Context, id string) error {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "DeleteRelationShip")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	tx := pr.db.Delete(&domain.Relationship{}, "id = ?", id)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return fmt.Errorf("no rows deleted")
	}

	return nil
}
