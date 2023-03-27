// Package domain holds all domain related code.
package domain

import (
	"time"

	"gorm.io/gorm"
)

type Person struct {
	ID          string         `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string         `json:"name,omitempty"`
	Parents     []Relationship `json:"parents,omitempty" gorm:"foreignKey:ParentID;references:ID;-"`
	Children    []Relationship `json:"children,omitempty" gorm:"foreignKey:ChildrenID;references:ID;-"`
	Siblings    []Relationship `json:"siblings,omitempty" gorm:"-"`
	Spouse      *Relationship  `json:"spouse,omitempty" gorm:"-"`
	BaconNumber int            `json:"bacon_number,omitempty" gorm:"-"`
	CreatedAt   *time.Time     `json:"createdAt"`
	UpdatedAt   *time.Time     `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type Relationship struct {
	ID        string         `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ParentID  string         `json:"parent"`
	ChildID   string         `json:"children"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// Cousins represents the cousin relationship between two persons
type Cousins struct {
	Person1 Person   `json:"person1"`
	Person2 Person   `json:"person2"`
	Shared  []Person `json:"shared"`
}

// BaconNumber represents the bacon's number between two persons
type BaconNumber struct {
	Person1   Person `json:"person1"`
	Person2   Person `json:"person2"`
	Number    int    `json:"number"`
	Path      []int  `json:"path"`
	Visited   []bool `json:"-"`
	Processed []bool `json:"-"`
}

type FamilyTreeNode struct {
	Person   Person
	Children []*FamilyTreeNode
	Parent   *FamilyTreeNode
}

type FamilyTree struct {
	Root *FamilyTreeNode
}

type FamilyMember struct {
	Name          string               `json:"name"`
	Relationships []FamilyRelationship `json:"relationships"`
}

type FamilyRelationship struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
}
