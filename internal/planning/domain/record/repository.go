package record

import (
	"context"
	"errors"
	"iter"
	"time"

	"github.com/google/uuid"
)

// ErrNoRecordFound is an error when a record is not found.
var ErrNoRecordFound = errors.New("record not found")

// Repository is a domain service interface for repository.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*IntakeRecord, error)
	GetByPlanID(ctx context.Context, planID uuid.UUID) ([]*IntakeRecord, error)
	Save(ctx context.Context, record *IntakeRecord) error
	UpdateByID(ctx context.Context, record *IntakeRecord) error
	SaveBulk(ctx context.Context, records []*IntakeRecord) error
	RecordsByTime(ctx context.Context, time time.Time) (iter.Seq[*IntakeRecord], error)
}
