// Package adapters implements all the necessary
// logic to talk to any external system.
package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/bhborges/family-tree-api/internal/app"
	"github.com/bhborges/family-tree-api/internal/domain"

	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// [application name]:[layer].
const _SegmentPrefix = "familytree:repository"

// PostgresRepository implements a PostgreSQL
// approach for family tree repository.
type PostgresRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewPostgresRepository provides a new instance of
// a PostgreSQL family tree repository.
func NewPostgresRepository(db *gorm.DB, log *zap.Logger) *PostgresRepository {
	return &PostgresRepository{db, log}
}

// ListPeople returns a list with all people registered.
func (pr *PostgresRepository) ListPeople(ctx context.Context) (
	[]*domain.Person, error,
) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "ListAllPeople")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	var p []*domain.Person

	tx := pr.db.WithContext(ctx)
	tx.Where("deleted_at IS NULL")

	err := tx.Find(&p).Error
	if err != nil {
		return nil, err
	}

	return p, err
}

// GetPerson returns a person registered.
// Filtered by ID.
func (pr *PostgresRepository) GetPersonByID(ctx context.Context, id string) (*domain.Person, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "GetPerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	var p domain.Person

	tx := pr.db.WithContext(ctx)
	err := tx.Where(&domain.Person{
		ID: id,
	}).First(&p).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, app.ErrPersonNotFound
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// CreatePerson create a new person.
func (pr *PostgresRepository) CreatePerson(ctx context.Context, dp domain.Person) (string, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "CreatePerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	p := domain.Person{
		Name: dp.Name,
	}

	tx := pr.db.Create(&p)
	if tx.Error != nil {
		return "", tx.Error
	}

	return p.ID, nil
}

// CreatePeople creates a new batch of people.
func (pr *PostgresRepository) CreatePeople(ctx context.Context, people []domain.Person) ([]string, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "CreatePeople")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	var ids []string
	for _, p := range people {
		id, err := pr.CreatePerson(ctx, p)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// UpdatePerson update a person.
func (pr *PostgresRepository) UpdatePerson(ctx context.Context, dp domain.Person) error {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "UpdatePerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	tx := pr.db.Model(&dp).Update("name", dp.Name)

	return tx.Error
}

// DeletePerson delete a person.
func (pr *PostgresRepository) DeletePerson(ctx context.Context, id string) error {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "DeletePerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	tx := pr.db.Delete(&domain.Person{ID: id})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

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

// BuildFamilyTree builds the family tree of a given person ID, with the person as the root node.
func (pr *PostgresRepository) BuildFamilyTree(ctx context.Context, id string) (*domain.FamilyTree, error) {
	q := pr.db.Debug().Raw(`
		WITH RECURSIVE family_tree AS (
			SELECT id, parent_id, child_id
			FROM relationships
			WHERE child_id = ?
			UNION ALL
			SELECT r.id, r.parent_id, r.child_id
			FROM relationships r
			JOIN family_tree ft ON r.child_id = ft.parent_id
		)
		SELECT DISTINCT p.name as name, COALESCE(p2.name, '') as parent
		FROM family_tree ft
		JOIN people p ON ft.parent_id = p.id
		LEFT JOIN relationships r ON r.child_id = p.id
		LEFT JOIN people p2 ON r.parent_id = p2.id
		UNION ALL
		SELECT DISTINCT p.name as name, COALESCE(p2.name, '') as parent
		FROM people p
		LEFT JOIN relationships r ON r.child_id = p.id
		LEFT JOIN people p2 ON r.parent_id = p2.id
		WHERE p.id = ?`, id, id)
	rows, err := q.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := make(map[string][]string)

	for rows.Next() {
		var name string
		var parent string
		err := rows.Scan(&name, &parent)
		if err != nil {
			return nil, err
		}
		if parent != "" {
			if _, ok := r[name]; !ok {
				r[name] = []string{parent}
			} else {
				r[name] = append(r[name], parent)
			}
		}
	}

	ms := make([]*domain.Member, 0)
	for name, parentList := range r {
		m := &domain.Member{Name: name}
		fr := make([]domain.FamilyRelationship, 0)
		for _, parent := range parentList {
			r := domain.FamilyRelationship{Name: parent, Relationship: "parent"}
			fr = append(fr, r)
		}
		m.Relationships = fr
		ms = append(ms, m)
	}

	return &domain.FamilyTree{Members: ms}, nil
}
