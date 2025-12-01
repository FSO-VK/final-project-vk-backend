package plan

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrNoPlanFound is an error when a plan is not found.
var ErrNoPlanFound = errors.New("plan not found")

// Repository is a domain service interface for repository.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Plan, error)
	UserPlans(ctx context.Context, userID uuid.UUID) ([]*Plan, error)
	Save(ctx context.Context, plan *Plan) error
}
