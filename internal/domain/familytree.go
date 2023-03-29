// Package domain holds all domain related code.
package domain

import (
	"time"

	"gorm.io/gorm"
)

// Person represents a person or member.
type Person struct {
	ID          string         `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string         `json:"name,omitempty"`
	Parents     []*Person      `json:"parents,omitempty" gorm:"many2many:relationships;ForeignKey:ID;References:id"`
	Children    []*Person      `json:"children,omitempty" gorm:"many2many:relationships;ForeignKey:ID;References:id"`
	Siblings    []*Person      `json:"siblings,omitempty" gorm:"-"`
	Spouse      *Person        `json:"spouse,omitempty" gorm:"-"`
	BaconNumber int            `json:"baconNumber,omitempty" gorm:"-"`
	CreatedAt   *time.Time     `json:"createdAt"`
	UpdatedAt   *time.Time     `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index;-"`
}

// Relationship represents a many-to-many relationship between two persons.
type Relationship struct {
	ID        string         `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ParentID  string         `json:"parent" gorm:"primaryKey"`
	ChildID   string         `json:"children" gorm:"primaryKey"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index;-"`
}

// Cousins represents the cousin relationship between two persons.
type Cousins struct {
	Person1 Person   `json:"person1"`
	Person2 Person   `json:"person2"`
	Shared  []Person `json:"shared"`
}

// BaconNumber represents the bacon's number between two persons.
type BaconNumber struct {
	Person1   Person `json:"person1"`
	Person2   Person `json:"person2"`
	Number    int    `json:"number"`
	Path      []int  `json:"path"`
	Visited   []bool `json:"-"`
	Processed []bool `json:"-"`
}

// FamilyTree represents a collection of family members.
type FamilyTree struct {
	Members []*Member `json:"members"`
}

// Member represents a family member.
type Member struct {
	Name          string               `json:"name"`
	Relationships []FamilyRelationship `json:"relationships"`
}

// FamilyRelationship represents the relationship between two family members.
type FamilyRelationship struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
}
