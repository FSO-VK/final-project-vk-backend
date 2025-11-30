// Package application is a package for application logic of the planning service.
package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// ErrPlanNotFound is an error when plan is not belongs to user.
var ErrPlanNotBelongToUser = errors.New("plan does not belong to user")

// GetPlan is an interface for getting a notification.
type GetPlan interface {
	Execute(
		ctx context.Context,
		cmd *GetPlanCommand,
	) (*GetPlanResponse, error)
}

// GetPlanService is a service for creating a subscription.
type GetPlanService struct {
	planningRepo plan.Repository
	validator    validator.Validator
}

// NewGetPlanService returns a new GetPlanService.
func NewGetPlanService(
	planningRepo plan.Repository,
	valid validator.Validator,
) *GetPlanService {
	return &GetPlanService{
		planningRepo: planningRepo,
		validator:    valid,
	}
}

// GetPlanCommand is a request to get a plan.
type GetPlanCommand struct {
	ID     string `validate:"required,uuid"`
	UserID string `validate:"required,uuid"`
}

// GetPlanResponse is a response to get a plan.
type GetPlanResponse struct {
	ID             string
	MedicationID   string
	UserID         string
	AmountValue    float64
	AmountUnit     string
	Condition      string
	StartDate      string
	EndDate        string
	RecurrenceRule []string
}

// Execute executes the GetPlan command.
func (s *GetPlanService) Execute(
	ctx context.Context,
	req *GetPlanCommand,
) (*GetPlanResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("request is not valid: %w", valErr)
	}
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}

	parsedUser, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}

	requestedPlan, err := s.planningRepo.GetByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	if requestedPlan.UserID() != parsedUser {
		return nil, ErrPlanNotBelongToUser
	}

	amountValue, amountUnit := requestedPlan.Dosage()

	response := &GetPlanResponse{
		ID:             requestedPlan.ID().String(),
		MedicationID:   requestedPlan.MedicationID().String(),
		UserID:         requestedPlan.UserID().String(),
		AmountValue:    amountValue,
		AmountUnit:     amountUnit,
		Condition:      requestedPlan.Condition(),
		StartDate:      requestedPlan.CourseStart().Format(time.RFC3339),
		EndDate:        requestedPlan.CourseEnd().Format(time.RFC3339),
		RecurrenceRule: requestedPlan.ScheduleIcal(),
	}
	return response, nil
}
