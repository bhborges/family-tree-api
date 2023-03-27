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

// MyFamilyTree extends the domain FamilyTree structure for adding helper methods.
type MyFamilyTree struct {
	*domain.FamilyTree
}

type Person struct {
	Name string
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
	// TODO: fix
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

// BuildFamilyTree builds the family tree of a given person ID, with the person as the root node.
func (pr *PostgresRepository) BuildFamilyTree(ctx context.Context, id string) (*domain.FamilyTree, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "ListFamilyTreeByPerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	p, err := pr.GetPersonByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get person with ID %s: %v", id, err)
	}

	// Create a new family tree with the person as the root node.
	familyTree := NewFamilyTree(*p)

	// Get all parents of the person.
	parents, err := pr.getParents(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get parents for person with ID %s: %v", id, err)
	}

	// For each parent, recursively build the family tree.
	for _, parent := range parents {
		parentNode, err := pr.BuildFamilyTree(ctx, parent.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to build family tree for parent with ID %s: %v", parent.ID, err)
		}

		myFamilyTree := &MyFamilyTree{familyTree}
		// Create a new FamilyMember instance with the parent's name.
		parentMember := domain.FamilyMember{Name: parentNode.Root.Person.Name}
		// Add the parent node as a child of the root node.
		myFamilyTree.AddChild(familyTree.Root, parentMember)
	}

	return familyTree, nil
}

// FindNodeByName finds a node in the family tree by name.
func (t *MyFamilyTree) FindNodeByName(name string, node *domain.FamilyTreeNode) (*domain.FamilyTreeNode, error) {
	if node.Person.Name == name {
		return node, nil
	}

	for _, childNode := range node.Children {
		foundNode, err := t.FindNodeByName(name, childNode)
		if err == nil {
			return foundNode, nil
		}
	}

	return nil, fmt.Errorf("person with name %s not found in family tree", name)
}

func (pr *PostgresRepository) getParents(personID string) ([]domain.Person, error) {
	var parents []domain.Person
	err := pr.db.Model(&domain.Person{}).Where("id in (select parent_id from relationships where child_id = ?)", personID).Find(&parents).Error
	if err != nil {
		return nil, err
	}
	return parents, nil
}

func (t *MyFamilyTree) AddChild(parentNode *domain.FamilyTreeNode, childPerson domain.FamilyMember) error {
	// create a new node with the child person data
	childNode := &domain.FamilyTreeNode{
		Person: domain.Person{
			Name: childPerson.Name,
		},
	}

	// check if parent node exists
	if parentNode == nil {
		return errors.New("parent node is nil")
	}

	// add child node to parent node's children
	parentNode.Children = append(parentNode.Children, childNode)

	return nil
}

func (m Person) ToFamilyMember() domain.FamilyMember {
	return domain.FamilyMember{
		Name: m.Name,
	}
}

// NewFamilyTree creates a new family tree with the given person as the root node.
func NewFamilyTree(person domain.Person) *domain.FamilyTree {
	return &domain.FamilyTree{
		Root: &domain.FamilyTreeNode{
			Person: person,
		},
	}
}
