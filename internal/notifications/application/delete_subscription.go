package application

import (
	"context"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
)

// DeleteSubscription is an interface for adding a notification.
type DeleteSubscription interface {
	Execute(
		ctx context.Context,
		cmd *DeleteSubscriptionCommand,
	) (*DeleteSubscriptionResponse, error)
}

// DeleteSubscriptionService is a service for getting public key.
type DeleteSubscriptionService struct {
	subscriptionsRepo subscriptions.Repository
	validator validator.Validator
}

// NewDeleteSubscriptionService returns a new DeleteSubscriptionService.
func NewDeleteSubscriptionService(
	subscriptionsRepo subscriptions.Repository,
	valid validator.Validator,
) *DeleteSubscriptionService {
	return &DeleteSubscriptionService{
		subscriptionsRepo: subscriptionsRepo,
		validator: valid,
	}
}

// DeleteSubscriptionCommand is a request to to get public key.
type DeleteSubscriptionCommand struct {
	// TODO
}

// DeleteSubscriptionResponse is a response to get public key.
type DeleteSubscriptionResponse struct {
	// TODO
}

// Execute executes the DeleteSubscription command.
func (s *DeleteSubscriptionService) Execute(
	ctx context.Context,
	req *DeleteSubscriptionCommand,
) (*DeleteSubscriptionResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}
	// TODO
	return nil, nil
}
