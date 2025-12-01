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
	cfg         ClientConfig
	recordsRepo record.Repository
	planRepo    plan.Repository
	ticker      TickerInterface
}

// NewGenerateRecordService creates a new GenerateRecordService.
func NewGenerateRecordService(
	cfg ClientConfig,
	recordsRepo record.Repository,
	planRepo plan.Repository,
	ticker TickerInterface,
) *GenerateRecordService {
	return &GenerateRecordService{
		cfg:         cfg,
		recordsRepo: recordsRepo,
		planRepo:    planRepo,
		ticker:      ticker,
	}
}

func (g *GenerateRecordService) GenerateRecord(ctx context.Context, planID uuid.UUID) error {
	p, err := g.planRepo.GetByID(ctx, planID)
	if err != nil {
		return err
	}

	records, err := p.GenerateIntakeRecords(
		time.Now(),
		time.Now().Truncate(24*time.Hour).Add(g.cfg.CreationShift),
	)
	if err != nil {
		return err
	}

	if err := g.recordsRepo.SaveBulk(ctx, records); err != nil {
		return err
	}

	return nil
}

func (g *GenerateRecordService) GenerateRecordsForDay(ctx context.Context) error {
	fmt.Println("GenerateRecordsForDay")
	seq, err := g.planRepo.IterActivePlans(ctx, g.cfg.BatchSize)
	if err != nil {
		return err
	}

	creationTime := time.Now().Truncate(24 * time.Hour).Add(g.cfg.CreationShift)
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

func (d *GenerateRecordService) Run(ctx context.Context) {
	d.ticker.Run(ctx, func(ctx context.Context) {
		_ = d.GenerateRecordsForDay(ctx)
	})
}
