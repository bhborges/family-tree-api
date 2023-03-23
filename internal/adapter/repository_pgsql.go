// Package adapters implements all the necessary
// logic to talk to any external system.
package adapter

import (
	"context"
	"fmt"

	"github.com/bhborges/family-tree-api/internal/domain"

	"github.com/lib/pq"
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

// GetPerson returns a person registered.
// Filtered by ID.
func (r *PostgresRepository) GetPersonByID(ctx context.Context, ID string) (*domain.Person, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "GetPerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	var p domain.Person

	tx := r.db.WithContext(ctx)
	tx.Where("deleted_at IS NULL")

	if len(ID) > 0 {
		tx = tx.Where("code = ANY(?)", pq.Array(ID))
	}

	err := tx.Find(&p).Error
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// CreatePerson create a new person.
func (r *PostgresRepository) CreatePerson(ctx context.Context, p *domain.Person) (string, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "CreatePerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	tx := r.db.Create(&p)
	// TODO: fix
	if tx.Error != nil {
		return "", tx.Error
	}

	return p.ID, nil
}

// UpdatePerson update a person.
func (r *PostgresRepository) UpdatePerson(ctx context.Context, p *domain.Person) (string, error) {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "UpdatePerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	tx := r.db.Save(&p)
	// TODO: fix
	if tx.Error != nil {
		return "", tx.Error
	}

	return p.ID, nil
}

// DeletePerson delete a person.
func (r *PostgresRepository) DeletePerson(ctx context.Context, p *domain.Person) error {
	trans := newrelic.FromContext(ctx)
	if trans != nil {
		segmentName := fmt.Sprintf("%s:%s", _SegmentPrefix, "DeletePerson")
		segment := trans.StartSegment(segmentName)

		defer segment.End()
	}

	tx := r.db.Delete(&p)
	// TODO: fix
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
