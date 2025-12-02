package generaterecord

import (
	"context"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/record"
	"github.com/google/uuid"
)

// GenerateRecordProvider is an interface for generating records.
type GenerateRecordProvider interface {
	GenerateRecord(planID uuid.UUID) error
	GenerateRecordsForDay() error
}

// GenerateRecordService implements GenerateRecordProvider.
type GenerateRecordService struct {
	creationShift time.Duration
	batchSize     int
	recordsRepo   record.Repository
	planRepo      plan.Repository
}

// NewGenerateRecordService creates a new GenerateRecordService.
func NewGenerateRecordService(
	creationShift time.Duration,
	batchSize int,
	recordsRepo record.Repository,
	planRepo plan.Repository,
) *GenerateRecordService {
	return &GenerateRecordService{
		creationShift: creationShift,
		batchSize:     batchSize,
		recordsRepo:   recordsRepo,
		planRepo:      planRepo,
	}
}

// GenerateRecord generates records for a specific plan.
func (g *GenerateRecordService) GenerateRecord(ctx context.Context, planID uuid.UUID) error {
	p, err := g.planRepo.GetByID(ctx, planID)
	if err != nil {
		return err
	}

	records, err := p.GenerateIntakeRecords(
		time.Now(),
		time.Now().Truncate(24*time.Hour).Add(g.creationShift),
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
func (g *GenerateRecordService) GenerateRecordsForDay(ctx context.Context) error {
	fmt.Println("111111")
	seq, err := g.planRepo.ActivePlans(ctx, g.batchSize)
	if err != nil {
		return err
	}

	creationTime := time.Now().Truncate(24 * time.Hour).Add(g.creationShift)
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
