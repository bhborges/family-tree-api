// Package adapters implements all the necessary
// logic to talk to any external system.
package adapter

import (
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
