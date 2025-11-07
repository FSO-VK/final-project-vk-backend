package record

import (
	"context"

	"github.com/google/uuid"
)

// Repository is a domain service interface for repository.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*IntakeRecord, error)
	GetByPlanID(ctx context.Context, planID uuid.UUID) ([]*IntakeRecord, error)
	Save(ctx context.Context, record IntakeRecord) error
	SaveBulk(ctx context.Context, records []*IntakeRecord) error
}
