package application

import (
	"context"
	"fmt"

	provider "github.com/FSO-VK/final-project-vk-backend/internal/notifications/application/notification_provider"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// SendNotification is an interface for adding a notification.
type SendNotification interface {
	Execute(
		ctx context.Context,
		cmd *SendNotificationCommand,
	) (*SendNotificationResponse, error)
}

// SendNotificationService is a service for sending a notification.
type SendNotificationService struct {
	subscriptionsRepo    subscriptions.Repository
	notificationProvider provider.NotificationProvider
	validator            validator.Validator
}

// NewSendNotificationService returns a new SendNotificationService.
func NewSendNotificationService(
	subscriptionsRepo subscriptions.Repository,
	notificationProvider provider.NotificationProvider,
	valid validator.Validator,
) *SendNotificationService {
	return &SendNotificationService{
		subscriptionsRepo:    subscriptionsRepo,
		notificationProvider: notificationProvider,
		validator:            valid,
	}
}

// SendNotificationCommand is a request to send a notification.
type SendNotificationCommand struct {
	UserID string
	Title  string
	Body   string
}

// SendNotificationResponse is a response to send a notification.
type SendNotificationResponse struct{}

// Execute executes the SendNotification command.
func (s *SendNotificationService) Execute(
	ctx context.Context,
	req *SendNotificationCommand,
) (*SendNotificationResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}

	parsedUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}

	subscriptions, err := s.subscriptionsRepo.GetSubscriptionsByUserID(ctx, parsedUserID)
	if err != nil {
		return nil, fmt.Errorf("there is no such subscription: %w", err)
	}
	for _, subscription := range subscriptions {
		notificationToSend := provider.NewNotification(
			uuid.New(),
			subscription.GetID(),
			parsedUserID,
			req.Title,
			req.Body,
		)
		err = s.notificationProvider.PushNotification(notificationToSend, subscription)
		if err != nil {
			return nil, fmt.Errorf("failed to send notification: %w", err)
		}
	}

	response := &SendNotificationResponse{}
	return response, nil
}
