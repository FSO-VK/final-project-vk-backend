package application

import (
	"context"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/notifications"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
)

// InteractWithNotification is an interface for adding a notification.
type InteractWithNotification interface {
	Execute(
		ctx context.Context,
		cmd *InteractWithNotificationCommand,
	) (*InteractWithNotificationResponse, error)
}

// InteractWithNotificationService is a service for getting public key.
type InteractWithNotificationService struct {
	subscriptionsRepo subscriptions.Repository
	notificationsRepo notifications.Repository
	validator validator.Validator
}

// NewInteractWithNotificationService returns a new InteractWithNotificationService.
func NewInteractWithNotificationService(
	subscriptionsRepo subscriptions.Repository,
	notificationsRepo notifications.Repository,
	valid validator.Validator,
) *InteractWithNotificationService {
	return &InteractWithNotificationService{
		subscriptionsRepo: subscriptionsRepo,
		notificationsRepo: notificationsRepo,
		validator: valid,
	}
}

// InteractWithNotificationCommand is a request to to get public key.
type InteractWithNotificationCommand struct {
	// TODO
}

// InteractWithNotificationResponse is a response to get public key.
type InteractWithNotificationResponse struct {
	// TODO
}

// Execute executes the InteractWithNotification command.
func (s *InteractWithNotificationService) Execute(
	ctx context.Context,
	req *InteractWithNotificationCommand,
) (*InteractWithNotificationResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}
	// TODO
	return nil, nil
}
