package category

import "context"

// Repository is a domain interface for category's repository.
type Repository interface {
	Create(ctx context.Context, category *Category) (*Category, error)
	GetByID(ctx context.Context, categoryID uint) (*Category, error)
}
