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

// FinishPlan is an interface for completing a notification.
type FinishPlan interface {
	Execute(
		ctx context.Context,
		cmd *FinishPlanCommand,
	) (*FinishPlanResponse, error)
}

// FinishPlanService is a service for creating a subscription.
type FinishPlanService struct {
	planningRepo plan.Repository
	validator    validator.Validator
}

// NewFinishPlanService returns a new FinishPlanService.
func NewFinishPlanService(
	planningRepo plan.Repository,
	valid validator.Validator,
) *FinishPlanService {
	return &FinishPlanService{
		planningRepo: planningRepo,
		validator:    valid,
	}
}

// FinishPlanCommand is a request to finish a plan.
type FinishPlanCommand struct {
	ID     string `validate:"required,uuid"`
	UserID string `validate:"required,uuid"`
}

// FinishPlanResponse is a response to finish a plan.
type FinishPlanResponse struct{}

// Execute executes the FinishPlan command.
func (s *FinishPlanService) Execute(
	ctx context.Context,
	req *FinishPlanCommand,
) (*FinishPlanResponse, error) {
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
		return nil, fmt.Errorf("failed to finish plan: %w", err)
	}

	if p.UserID() != parsedUser {
		return nil, ErrNotOwner
	}

	newPlan, err := p.Deactivate()
	if err != nil {
		return nil, fmt.Errorf("failed to finish plan: %w", err)
	}

	err = s.planningRepo.UpdatePlan(ctx, newPlan)
	if err != nil {
		return nil, fmt.Errorf("failed to finish plan: %w", err)
	}

	return &FinishPlanResponse{}, nil
}
