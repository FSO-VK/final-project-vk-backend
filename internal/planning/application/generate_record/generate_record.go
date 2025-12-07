package generaterecord

import (
	"context"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/record"
	"github.com/google/uuid"
)

// GenerateRecord is an interface for generating records.
type GenerateRecord interface {
	GenerateRecord(planID uuid.UUID) error
	GenerateRecordsForDay() error
}

// GenerateRecordService implements GenerateRecord.
type GenerateRecordService struct {
	recordsRepo record.Repository
	planRepo    plan.Repository
}

// NewGenerateRecordService creates a new GenerateRecordService.
func NewGenerateRecordService(
	recordsRepo record.Repository,
	planRepo plan.Repository,
) *GenerateRecordService {
	return &GenerateRecordService{
		recordsRepo: recordsRepo,
		planRepo:    planRepo,
	}
}

// GenerateRecord generates records for a specific plan.
func (g *GenerateRecordService) GenerateRecord(
	ctx context.Context,
	planID uuid.UUID,
	creationShift time.Duration,
) error {
	p, err := g.planRepo.GetByID(ctx, planID)
	if err != nil {
		return err
	}

	records, err := p.GenerateIntakeRecords(
		time.Now(),
		time.Now().Truncate(24*time.Hour).Add(creationShift),
	)
	if err != nil {
		return err
	}

	if err := g.recordsRepo.SaveBulk(ctx, records); err != nil {
		return err
	}

	return nil
}

// GenerateRecordsForDay generates records for all active plans.
func (g *GenerateRecordService) GenerateRecordsForDay(
	ctx context.Context,
	batchSize int,
	creationShift time.Duration,
) error {
	seq, err := g.planRepo.ActivePlans(ctx, batchSize)
	if err != nil {
		return err
	}

	creationTime := time.Now().Truncate(24 * time.Hour).Add(creationShift)
	now := time.Now()

	for p := range seq {
		records, err := p.GenerateIntakeRecords(now, creationTime)
		if err != nil {
			return err
		}

		if err := g.recordsRepo.SaveBulk(ctx, records); err != nil {
			return err
		}
	}

	return nil
}
