// Package category is a package for categories of medications
package category

// Category represents a category.
type Category struct {
	ID   uint
	Name string
}

// NewCategory creates a new category.
func NewCategory(name string) *Category {
	return &Category{
		ID:   0, // could not be nill, but will be set later
		Name: name,
	}
}
