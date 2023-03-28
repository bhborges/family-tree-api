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
	"gorm.io/gorm"
)

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
