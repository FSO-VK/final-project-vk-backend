package application

import (
	"context"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/domain/subscriptions"
)

// CreateSubscription is an interface for adding a notification.
type CreateSubscription interface {
	Execute(
		ctx context.Context,
		cmd *CreateSubscriptionCommand,
	) (*CreateSubscriptionResponse, error)
}

// CreateSubscriptionService is a service for getting public key.
type CreateSubscriptionService struct {
	subscriptionsRepo subscriptions.Repository
	validator validator.Validator
}

// NewCreateSubscriptionService returns a new CreateSubscriptionService.
func NewCreateSubscriptionService(
	subscriptionsRepo subscriptions.Repository,
	valid validator.Validator,
) *CreateSubscriptionService {
	return &CreateSubscriptionService{
		subscriptionsRepo: subscriptionsRepo,
		validator: valid,
	}
}

// CreateSubscriptionCommand is a request to to get public key.
type CreateSubscriptionCommand struct {
	// TODO
}

// CreateSubscriptionResponse is a response to get public key.
type CreateSubscriptionResponse struct {
	// TODO
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
	// TODO
	return nil, nil
}
