package application

import (
	"context"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/application/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/application/notification"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/record"
)

// IntakeNotificationGenerator is an interface for generating notifications for intake.
type IntakeNotificationGenerator interface {
	GenerateIntakeNotifications(ctx context.Context) error
}

// IntakeNotificationService implements IntakeNotification.
type IntakeNotificationService struct {
	recordsRepo          record.Repository
	planRepo             plan.Repository
	notificationProvider notification.NotificationService
	medicationProvider   medication.MedicationService
}

// NewIntakeNotificationService creates a new IntakeNotificationService.
func NewIntakeNotificationService(
	recordsRepo record.Repository,
	planRepo plan.Repository,
	notificationProvider notification.NotificationService,
	medicationProvider medication.MedicationService,
) *IntakeNotificationService {
	return &IntakeNotificationService{
		recordsRepo:          recordsRepo,
		planRepo:             planRepo,
		notificationProvider: notificationProvider,
		medicationProvider:   medicationProvider,
	}
}

// GenerateIntakeNotifications generates notifications for intake.
func (g *IntakeNotificationService) GenerateIntakeNotifications(
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
		medicationName, err := g.medicationProvider.MedicationName(p.MedicationID(), p.UserID())
		if err != nil {
			continue
		}
		info := notification.NotificationInfo{
			UserID: p.UserID(),
			Title:  "Время принять " + medicationName,
			Body:   p.Condition(),
		}
		if err := g.notificationProvider.SendNotification(ctx, info); err != nil {
			return err
		}
	}
	return nil
}
