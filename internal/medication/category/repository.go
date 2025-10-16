package category

import "context"

// RepositoryForCategory is a repository for categories and provides methods to create and get by id categories.
type RepositoryForCategory interface {
	Create(ctx context.Context, category *Category) (*Category, error)
	GetByID(ctx context.Context, categoryID uint) (*Category, error)
}
