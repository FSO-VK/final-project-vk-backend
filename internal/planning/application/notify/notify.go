package notify

import (
	"context"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/application/notification"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/record"
	client "github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/notification_client"
)

// IntakeNotification is an interface for generating notifications for intake.
type IntakeNotification interface {
	GenerateNotifications(ctx context.Context) error
}

// IntakeNotificationService implements IntakeNotification.
type IntakeNotificationService struct {
	recordsRepo          record.Repository
	planRepo             plan.Repository
	notificationProvider client.NotificationClient
}

// NewIntakeNotificationService creates a new IntakeNotificationService.
func NewIntakeNotificationService(
	recordsRepo record.Repository,
	planRepo plan.Repository,
	notificationProvider client.NotificationClient,
) *IntakeNotificationService {
	return &IntakeNotificationService{
		recordsRepo:          recordsRepo,
		planRepo:             planRepo,
		notificationProvider: notificationProvider,
	}
}

// GenerateRecordsForDay generates records for all active plans.
func (g *IntakeNotificationService) GenerateNotifications(
	ctx context.Context,
) error {
	records, err := g.recordsRepo.RecordsByTime(ctx, time.Now())
	if err != nil {
		return err
	}

	for r := range records {
		p, err := g.planRepo.GetByID(ctx, r.PlanID())
		if err != nil {
			continue
		}
		info := notification.NotificationInfo{
			UserID: p.UserID(),
			Title:  "Время принять таблетки", // + p.Name()
			Body:   p.Condition(),
		}
		if err := g.notificationProvider.SendNotification(ctx, info); err != nil {
			return err
		}
	}
	return nil
}
