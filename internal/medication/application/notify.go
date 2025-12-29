package application

import (
	"context"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application/notification"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medbox"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
)

// ExpirationNotificationGenerator is an interface for generating notifications for expiration date of medication.
type ExpirationNotificationGenerator interface {
	GenerateExpirationNotifications(ctx context.Context) error
}

// ExpirationNotificationService implements ExpirationNotification.
type ExpirationNotificationService struct {
	medicationRepo       medication.Repository
	medBoxRepo           medbox.Repository
	notificationProvider notification.NotificationService
}

// NewExpirationNotificationService creates a new ExpirationNotificationService.
func NewExpirationNotificationService(
	medicationRepo medication.Repository,
	medBoxRepo medbox.Repository,
	notificationProvider notification.NotificationService,
) *ExpirationNotificationService {
	return &ExpirationNotificationService{
		medicationRepo:       medicationRepo,
		medBoxRepo:           medBoxRepo,
		notificationProvider: notificationProvider,
	}
}

// GenerateExpirationNotifications generates notifications for expiration date of medication.
func (g *ExpirationNotificationService) GenerateExpirationNotifications(
	ctx context.Context,
	timeDelta time.Duration,
) error {
	medications, err := g.medicationRepo.MedicationByExpiration(ctx, timeDelta)
	if err != nil {
		return err
	}

	for m := range medications {
		userID, err := g.medBoxRepo.GetUserByMedicationID(ctx, m.GetID())
		if err != nil {
			continue
		}
		info := notification.NotificationInfo{
			UserID: userID,
			Title:  "Истекает срок годности " + string(m.GetInternationalName()),
			Body: "Срок годности препарата " + string(m.GetInternationalName()) +
				" истекает " + m.GetExpirationDate().Format("02.01.2006"),
		}
		if err := g.notificationProvider.SendNotification(ctx, info); err != nil {
			return err
		}
	}
	return nil
}
