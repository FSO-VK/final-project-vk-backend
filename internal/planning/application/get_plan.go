// Package application is a package for application logic of the planning service.
package application

import (
	"context"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

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
	Status         string
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

	requestedPlan, err := s.planningRepo.GetByID(ctx, parsedID)
	if err != nil {
		return nil, ErrNoPlan
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
		Status:         requestedPlan.Status().String(),
		StartDate:      requestedPlan.CourseStart().Format(time.RFC3339),
		EndDate:        requestedPlan.CourseEnd().Format(time.RFC3339),
		RecurrenceRule: requestedPlan.ScheduleIcal(),
	}
	return response, nil
}
