// Package domain holds all domain related code.
package domain

// Person represents a person.
type Person struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Relationships represents with person.
type Relationships struct {
	ID       string  `json:"id,omitempty"`
	Parent   *Person `json:"parente"`
	Children *Person `json:"children"`
}
