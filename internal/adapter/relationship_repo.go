package adapter

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/bhborges/family-tree-api/internal/app"
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

	rs, err := pr.listRelationshipsByParentAndChildIDs(dr.ParentID, dr.ChildID)
	if err != nil {
		return "", err
	}

	if pr.checkIncestuousOffspring(rs) {
		return "", app.ErrIncestuousOffspring
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
		return app.ErrNoRowsUpdated
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
		return app.ErrNoRowsInserted
	}

	return nil
}

// listParentsByID lists all parents that have the given parentID.
func (pr *PostgresRepository) listParentsByID(parentID string) ([]*domain.Relationship, error) {
	relationships := []*domain.Relationship{}
	tx := pr.db.Debug().Where("child_id = ?", parentID).Find(&relationships)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return relationships, nil
}

// listRelationshipsByParentAndChildIDs lists all relationships that match the given parentID and childID.
func (pr *PostgresRepository) listRelationshipsByParentAndChildIDs(parentID, childID string) ([]*domain.Relationship, error) {
	rs := []*domain.Relationship{}

	// get all relationships where the given parent or child is involved
	tx := pr.db.Debug().Where(
		"(parent_id = ? OR child_id = ?) OR (parent_id = ? OR child_id = ?)", parentID, childID, childID, parentID,
	).Find(&rs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// check for additional relationships involving the parents
	rsParents := []*domain.Relationship{}
	tx = pr.db.Debug().Where(
		"parent_id IN (?) OR child_id IN (?)", parentID, parentID,
	).Find(&rsParents)

	if tx.Error != nil {
		return nil, tx.Error
	}

	rdChilds := []*domain.Relationship{}

	for _, r := range rsParents {
		if r.ParentID != parentID {
			// get all relationships where the parent's other child is involved
			tx = pr.db.Debug().Where(
				"(parent_id = ? AND child_id != ?) OR (parent_id = ? AND child_id != ?)", r.ParentID, childID, childID, r.ParentID,
			).Find(&rdChilds)
			if tx.Error != nil {
				return nil, tx.Error
			}

			rs = append(rs, rdChilds...)
		}
	}

	return rs, nil
}

// fix:
// CheckIncestuousOffspring checks if there are any incestuous relationships in a given slice.
func (pr *PostgresRepository) checkIncestuousOffspring(relationships []*domain.Relationship) bool {
	for i, r1 := range relationships {
		for _, r2 := range relationships[i+1:] {
			if (r1.ParentID == r2.ParentID && r1.ChildID == r2.ChildID) || (r1.ParentID == r2.ChildID && r1.ChildID == r2.ParentID) {
				return true
			} else {
				p1, err := pr.listParentsByID(r1.ParentID)
				if err != nil {
					pr.log.Error(fmt.Sprintf("Error while getting parents for ID: %s", r1.ParentID), zap.Error(err))

					return false
				}
				p2, err := pr.listParentsByID(r2.ParentID)
				if err != nil {
					pr.log.Error(fmt.Sprintf("Error while getting parents for ID: %s", r2.ParentID), zap.Error(err))

					return false
				}
				for _, pp1 := range p1 {
					for _, pp2 := range p2 {
						if pp1.ID == pp2.ID {
							return true
						}
					}
				}
			}
		}
	}

	return false
}
