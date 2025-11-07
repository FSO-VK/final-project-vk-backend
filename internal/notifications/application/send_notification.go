package application

import (
	"context"
	"fmt"
	"time"

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
	SubscriptionID string
	UserID         string
	Title          string
	Body           string
	SendAt         string
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
	parsedSubscriptionID, parsedUserID, parsedSendAt, err := purse(
		req.SubscriptionID,
		req.UserID,
		req.SendAt)
	if err != nil {
		return nil, err
	}

	subscription, err := s.subscriptionsRepo.GetSubscriptionByID(ctx, parsedSubscriptionID)
	if err != nil {
		return nil, fmt.Errorf("there is no such subscription: %w", err)
	}
	if subscription.GetUserID() != parsedUserID {
		return nil, fmt.Errorf("this is not your subscription: %w", err)
	}

	notificationToSend := provider.NewNotification(
		uuid.New(),
		subscription.GetID(),
		parsedUserID,
		req.Title,
		req.Body,
		parsedSendAt,
	)
	err = s.notificationProvider.PushNotification(notificationToSend, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}

	response := &SendNotificationResponse{}
	return response, nil
}

func purse(
	subscriptionID string,
	userID string,
	sendAt string,
) (uuid.UUID, uuid.UUID, time.Time, error) {
	parsedSubscriptionID, err := uuid.Parse(subscriptionID)
	if err != nil {
		return uuid.Nil, uuid.Nil, time.Time{}, fmt.Errorf(
			"subscription id is invalid: %w: %w",
			ErrDeleteInvalidUUIDFormat,
			err,
		)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, uuid.Nil, time.Time{}, fmt.Errorf(
			"user id is invalid: %w: %w",
			ErrDeleteInvalidUUIDFormat,
			err)
	}

	parsedSendAt, err := time.Parse(time.RFC3339, sendAt)
	if err != nil {
		return uuid.Nil, uuid.Nil, time.Time{}, fmt.Errorf(
			"invalid sendAt format: %w",
			err)
	}

	return parsedSubscriptionID, parsedUserID, parsedSendAt, nil
}
