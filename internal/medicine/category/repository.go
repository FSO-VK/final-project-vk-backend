package category

import "context"

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) (*Category, error)
	GetByID(ctx context.Context, categoryID uint) (*Category, error)
}