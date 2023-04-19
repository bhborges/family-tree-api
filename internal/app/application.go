// Package app holds the business logic.
package app

import (
	"context"

	"github.com/bhborges/family-tree-api/internal/domain"

	"go.uber.org/zap"
)

// [application name]:[layer].
const _SegmentPrefix = "familytree:application"

// Application definition.
type Application struct {
	repository Repository
	log        *zap.Logger
}

// Repository specifies the signature of a person repository.
type Repository interface {
	ListPeople(context.Context) ([]*domain.Person, error)
	ListRelationships(context.Context) ([]*domain.Relationship, error)
	GetPersonByID(context.Context, string) (*domain.Person, error)
	CreatePerson(context.Context, domain.Person) (string, error)
	CreatePeople(context.Context, []domain.Person) ([]string, error)
	UpdatePerson(context.Context, domain.Person) error
	DeletePerson(context.Context, string) error
	CreateRelationship(context.Context, domain.Relationship) (string, error)
	UpdateRelationship(context.Context, *domain.Relationship) error
	DeleteRelationship(context.Context, string) error
	BuildFamilyTree(context.Context, string) (*domain.FamilyTree, error)
}

// NewApplication initializes an instance of a person Application.
func NewApplication(repository Repository, log *zap.Logger) *Application {
	return &Application{repository, log}
}
