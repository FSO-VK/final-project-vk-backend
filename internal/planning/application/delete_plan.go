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

// CompletePlan is an interface for completing a notification.
type CompletePlan interface {
	Execute(
		ctx context.Context,
		cmd *CompletePlanCommand,
	) (*CompletePlanResponse, error)
}

// CompletePlanService is a service for creating a subscription.
type CompletePlanService struct {
	planningRepo plan.Repository
	validator    validator.Validator
}

// NewCompletePlanService returns a new CompletePlanService.
func NewCompletePlanService(
	planningRepo plan.Repository,
	valid validator.Validator,
) *CompletePlanService {
	return &CompletePlanService{
		planningRepo: planningRepo,
		validator:    valid,
	}
}

// CompletePlanCommand is a request to complete a plan.
type CompletePlanCommand struct {
	ID     string `validate:"required,uuid"`
	UserID string `validate:"required,uuid"`
}

// CompletePlanResponse is a response to complete a plan.
type CompletePlanResponse struct{}

// Execute executes the CompletePlan command.
func (s *CompletePlanService) Execute(
	ctx context.Context,
	req *CompletePlanCommand,
) (*CompletePlanResponse, error) {
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
		return nil, fmt.Errorf("failed to complete plan: %w", err)
	}

	if p.UserID() != parsedUser {
		return nil, ErrNotOwner
	}

	newPlan, err := p.Deactivate()
	if err != nil {
		return nil, fmt.Errorf("failed to complete plan: %w", err)
	}

	err = s.planningRepo.UpdatePlan(ctx, newPlan)
	if err != nil {
		return nil, fmt.Errorf("failed to complete plan: %w", err)
	}

	return &CompletePlanResponse{}, nil
}
