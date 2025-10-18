package category

import "context"

// RepositoryForCategory is a domain interface for category's repository.
type RepositoryForCategory interface {
	Create(ctx context.Context, category *Category) (*Category, error)
	GetByID(ctx context.Context, categoryID uint) (*Category, error)
}
