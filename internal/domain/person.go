// Package domain holds all domain related code.
package domain

import (
	"time"

	"gorm.io/gorm"
)

// Person represents a person.
type Person struct {
	ID   string `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name string `json:"name,omitempty"`
	CreatedAt   *time.Time        `json:"createdAt"`
	UpdatedAt   *time.Time        `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt    `json:"deletedAt" gorm:"index"`
}

// Relationships represents with people.
type Relationships struct {
	ID       string  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Parent   *Person `json:"parente"`
	Children *Person `json:"children"`
	CreatedAt   *time.Time        `json:"createdAt"`
	UpdatedAt   *time.Time        `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt    `json:"deletedAt" gorm:"index"`
}
