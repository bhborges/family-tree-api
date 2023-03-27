// Package domain holds all domain related code.
package model

import (
	"time"

	"gorm.io/gorm"
)

// Person represents a person.
type Person struct {
	ID        string         `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string         `json:"name,omitempty"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	// Parents   []Relationship `json:"parents" gorm:"foreignKey:ParentID;references:ID"`
	// Childrens []Relationship `json:"childrens" gorm:"foreignKey:ChildrenID;references:ID"`
}

// Relationship represents a kinship relationship between people.
type Relationship struct {
	ID        string         `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ParentID  string         `json:"parent"`
	ChildID   string         `json:"children"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
