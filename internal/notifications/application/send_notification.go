package application

import (
	"context"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/notifications"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
)

// SendNotification is an interface for adding a notification.
type SendNotification interface {
	Execute(
		ctx context.Context,
		cmd *SendNotificationCommand,
	) (*SendNotificationResponse, error)
}

// SendNotificationService is a service for getting public key.
type SendNotificationService struct {
	subscriptionsRepo subscriptions.Repository
	notificationsRepo notifications.Repository
	validator validator.Validator
}

// NewSendNotificationService returns a new SendNotificationService.
func NewSendNotificationService(
	subscriptionsRepo subscriptions.Repository,
	notificationsRepo notifications.Repository,
	valid validator.Validator,
) *SendNotificationService {
	return &SendNotificationService{
		subscriptionsRepo: subscriptionsRepo,
		notificationsRepo: notificationsRepo,
		validator: valid,
	}
}

// SendNotificationCommand is a request to to get public key.
type SendNotificationCommand struct {
	// TODO
}

// SendNotificationResponse is a response to get public key.
type SendNotificationResponse struct {
	// TODO
}

// Execute executes the SendNotification command.
func (s *SendNotificationService) Execute(
	ctx context.Context,
	req *SendNotificationCommand,
) (*SendNotificationResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}
	// TODO
	return nil, nil
}
