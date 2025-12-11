// Package application is a package for application logic of the planning service.
package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// ErrNotOwner is an error when user is not owner of the plan.
var ErrNotOwner = errors.New("user is not owner of the plan")

// DeletePlan is an interface for deleting a notification.
type DeletePlan interface {
	Execute(
		ctx context.Context,
		cmd *DeletePlanCommand,
	) (*DeletePlanResponse, error)
}

// DeletePlanService is a service for creating a subscription.
type DeletePlanService struct {
	planningRepo plan.Repository
	validator    validator.Validator
}

// NewDeletePlanService returns a new DeletePlanService.
func NewDeletePlanService(
	planningRepo plan.Repository,
	valid validator.Validator,
) *DeletePlanService {
	return &DeletePlanService{
		planningRepo: planningRepo,
		validator:    valid,
	}
}

// DeletePlanCommand is a request to delete a plan.
type DeletePlanCommand struct {
	ID     string `validate:"required,uuid"`
	UserID string `validate:"required,uuid"`
}

// DeletePlanResponse is a response to delete a plan.
type DeletePlanResponse struct{}

// Execute executes the DeletePlan command.
func (s *DeletePlanService) Execute(
	ctx context.Context,
	req *DeletePlanCommand,
) (*DeletePlanResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, ErrValidationFail
	}

	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, ErrValidationFail
	}
	parsedUser, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, ErrValidationFail
	}

	p, err := s.planningRepo.GetByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete plan: %w", err)
	}

	if p.UserID() != parsedUser {
		return nil, ErrNotOwner
	}

	newPlan, err := p.Deactivate()
	if err != nil {
		return nil, fmt.Errorf("failed to delete plan: %w", err)
	}

	err = s.planningRepo.UpdatePlan(ctx, newPlan)
	if err != nil {
		return nil, fmt.Errorf("failed to delete plan: %w", err)
	}

	return &DeletePlanResponse{}, nil
}
