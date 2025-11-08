package application

import (
	"context"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// CreateSubscription is an interface for adding a notification.
type CreateSubscription interface {
	Execute(
		ctx context.Context,
		cmd *CreateSubscriptionCommand,
	) (*CreateSubscriptionResponse, error)
}

// CreateSubscriptionService is a service for creating a subscription.
type CreateSubscriptionService struct {
	subscriptionsRepo subscriptions.Repository
	validator         validator.Validator
}

// NewCreateSubscriptionService returns a new CreateSubscriptionService.
func NewCreateSubscriptionService(
	subscriptionsRepo subscriptions.Repository,
	valid validator.Validator,
) *CreateSubscriptionService {
	return &CreateSubscriptionService{
		subscriptionsRepo: subscriptionsRepo,
		validator:         valid,
	}
}

// CreateSubscriptionCommand is a request to create subscription.
type CreateSubscriptionCommand struct {
	UserID    string
	SendInfo  SendInfo
	UserAgent string
}

// SendInfo is unique info for sending subscriptions.
type SendInfo struct {
	Endpoint string
	Keys     Keys
}

// Keys is unique keys for encryption.
type Keys struct {
	P256dh string
	Auth   string
}

// CreateSubscriptionResponse is a response to create subscription.
type CreateSubscriptionResponse struct {
	ID        string
	UserID    string
	SendInfo  SendInfo
	UserAgent string
	IsActive  bool
}

// Execute executes the CreateSubscription command.
func (s *CreateSubscriptionService) Execute(
	ctx context.Context,
	req *CreateSubscriptionCommand,
) (*CreateSubscriptionResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}
	parsedUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDeleteInvalidUUIDFormat, err)
	}
	subscription := subscriptions.NewSubscription(
		parsedUUID,
		req.SendInfo.Endpoint,
		req.SendInfo.Keys.P256dh,
		req.SendInfo.Keys.Auth,
		req.UserAgent,
	)

	err = s.subscriptionsRepo.SetSubscription(ctx, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to set subscription: %w", err)
	}
	subscriptionInBase, err := s.subscriptionsRepo.GetSubscriptionByID(ctx, subscription.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	response := &CreateSubscriptionResponse{
		ID:     subscriptionInBase.GetID().String(),
		UserID: subscriptionInBase.GetUserID().String(),
		SendInfo: SendInfo{
			Endpoint: subscriptionInBase.GetSendInfo().Endpoint,
			Keys: Keys{
				P256dh: subscriptionInBase.GetSendInfo().Keys.P256dh,
				Auth:   subscriptionInBase.GetSendInfo().Keys.Auth,
			},
		},
		UserAgent: subscription.GetUserAgent(),
		IsActive:  subscriptionInBase.GetIsActive(),
	}

	return response, nil
}
