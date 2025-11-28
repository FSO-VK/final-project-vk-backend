package plan

import (
	"context"

	"github.com/google/uuid"
)

// Repository is a domain service interface for repository.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Plan, error)
	UserPlans(ctx context.Context, userID uuid.UUID) (*[]Plan, error)
	Save(ctx context.Context, plan Plan) error
}
