package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// ErrDeleteInvalidUUIDFormat represents an error when the uuid is invalid.
var ErrDeleteInvalidUUIDFormat = errors.New("invalid UUID format")

// DeleteSubscription is an interface for adding a notification.
type DeleteSubscription interface {
	Execute(
		ctx context.Context,
		cmd *DeleteSubscriptionCommand,
	) (*DeleteSubscriptionResponse, error)
}

// DeleteSubscriptionService is a service for deleting a subscription.
type DeleteSubscriptionService struct {
	subscriptionsRepo subscriptions.Repository
	validator         validator.Validator
}

// NewDeleteSubscriptionService returns a new DeleteSubscriptionService.
func NewDeleteSubscriptionService(
	subscriptionsRepo subscriptions.Repository,
	valid validator.Validator,
) *DeleteSubscriptionService {
	return &DeleteSubscriptionService{
		subscriptionsRepo: subscriptionsRepo,
		validator:         valid,
	}
}

// DeleteSubscriptionCommand is a request to delete a subscription.
type DeleteSubscriptionCommand struct {
	ID string
}

// DeleteSubscriptionResponse is a response to delete a subscription.
type DeleteSubscriptionResponse struct{}

// Execute executes the DeleteSubscription command.
func (s *DeleteSubscriptionService) Execute(
	ctx context.Context,
	req *DeleteSubscriptionCommand,
) (*DeleteSubscriptionResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}
	parsedUUID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDeleteInvalidUUIDFormat, err)
	}
	subscription, err := s.subscriptionsRepo.GetSubscriptionByID(ctx, parsedUUID)
	if err != nil {
		return nil, fmt.Errorf("there is no such subscription: %w", err)
	}
	subscription.SetIsActive(false)
	err = s.subscriptionsRepo.SetSubscription(ctx, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to set subscription: %w", err)
	}
	response := &DeleteSubscriptionResponse{}
	return response, nil
}
