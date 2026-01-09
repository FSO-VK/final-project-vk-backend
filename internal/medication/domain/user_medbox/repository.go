package usermedbox

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, mb *UserMedbox) error
	GetOrCreate(ctx context.Context, userID uuid.UUID) (*UserMedbox, error)
}
