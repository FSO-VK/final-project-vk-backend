// Package category is a package for domain logic of medication categories.
package category

// Category is an entity of domain.
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
