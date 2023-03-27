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
	var rootPerson domain.Person

	if err := pr.db.Debug().
		Joins("JOIN relationships ON relationships.child_id = people.id").
		Preload("Parents").
		Preload("Children").
		Where("relationships.parent_id = ?", id).
		Find(&rootPerson.Children).
		Error; err != nil {
		return nil, err
	}

	tree := &domain.FamilyTree{
		Members: []*domain.Member{},
	}

	peopleMap := make(map[string]*domain.Person)

	peopleMap[rootPerson.ID] = &rootPerson

	childrenQueue := rootPerson.Children

	for len(childrenQueue) > 0 {
		child := childrenQueue[0]
		childrenQueue = childrenQueue[1:]

		childPerson := &domain.Person{}
		if err := pr.db.Preload("Parents").Preload("Children").First(childPerson, child.Children).Error; err != nil {
			continue
		}

		peopleMap[childPerson.ID] = childPerson

		relationships := buildFamilyRelationships(childPerson, peopleMap)

		member := &domain.Member{
			Name:          childPerson.Name,
			Relationships: relationships,
		}
		tree.Members = append(tree.Members, member)

		childrenQueue = append(childrenQueue, childPerson.Children...)
	}

	return tree, nil
}

// buildFamilyRelationships creates family relationships for a given person and returns them as a slice.
// It takes a map of people with person IDs as keys for efficient lookups.
func buildFamilyRelationships(person *domain.Person, peopleMap map[string]*domain.Person) []domain.FamilyRelationship {
	var relationships []domain.FamilyRelationship

	for _, parent := range person.Parents {
		parentPerson := peopleMap[parent.ID]
		if parentPerson != nil {
			relationship := domain.FamilyRelationship{
				Name:         parentPerson.Name,
				Relationship: "Parent",
			}
			relationships = append(relationships, relationship)
		}
	}

	for _, child := range person.Children {
		childPerson := peopleMap[child.ID]
		if childPerson != nil {
			relationship := domain.FamilyRelationship{
				Name:         childPerson.Name,
				Relationship: "Child",
			}
			relationships = append(relationships, relationship)
		}
	}

	return relationships
}
