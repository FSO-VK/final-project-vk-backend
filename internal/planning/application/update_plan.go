// Package application is a package for application logic of the planning service.
package application

import (
	"context"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
)

// UpdatePlan is an interface for adding a notification.
type UpdatePlan interface {
	Execute(
		ctx context.Context,
		cmd *UpdatePlanCommand,
	) (*UpdatePlanResponse, error)
}

// UpdatePlanService is a service for creating a subscription.
type UpdatePlanService struct {
	planningRepo plan.Repository
	validator    validator.Validator
}

// NewUpdatePlanService returns a new UpdatePlanService.
func NewUpdatePlanService(
	planningRepo plan.Repository,
	valid validator.Validator,
) *UpdatePlanService {
	return &UpdatePlanService{
		planningRepo: planningRepo,
		validator:    valid,
	}
}

// UpdatePlanCommand is a request to add a plan.
type UpdatePlanCommand struct {
	ID             string   `validate:"required,uuid"`
	MedicationID   string   `validate:"required,uuid"`
	UserID         string   `validate:"required,uuid"`
	AmountValue    float64  `validate:"required,gte=0"`
	AmountUnit     string   `validate:"required"`
	Condition      string   `validate:"omitempty,max=300"`
	StartDate      string   `validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate        string   `validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Duration       string   `validate:"required"`
	RecurrenceRule []string `validate:"required"`
}

// UpdatePlanResponse is a response to add a plan.
type UpdatePlanResponse struct {
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

// Execute executes the UpdatePlan command.
func (s *UpdatePlanService) Execute(
	ctx context.Context,
	req *UpdatePlanCommand,
) (*UpdatePlanResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, ErrValidationFail
	}

	draftPlan := &PlanDraft{
		ID:             req.ID,
		MedicationID:   req.MedicationID,
		UserID:         req.UserID,
		AmountValue:    req.AmountValue,
		AmountUnit:     req.AmountUnit,
		Condition:      req.Condition,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Duration:       req.Duration,
		RecurrenceRule: req.RecurrenceRule,
	}

	newPlan, err := createPlan(draftPlan)
	if err != nil {
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}

	err = s.planningRepo.Save(ctx, newPlan)
	if err != nil {
		return nil, fmt.Errorf("failed to save plan: %w", err)
	}
	amountValue, amountUnit := newPlan.Dosage()

	response := &UpdatePlanResponse{
		ID:             newPlan.ID().String(),
		MedicationID:   newPlan.MedicationID().String(),
		UserID:         newPlan.UserID().String(),
		AmountValue:    amountValue,
		AmountUnit:     amountUnit,
		Condition:      newPlan.Condition(),
		StartDate:      newPlan.CourseStart().Format(time.RFC3339),
		EndDate:        newPlan.CourseEnd().Format(time.RFC3339),
		RecurrenceRule: newPlan.ScheduleIcal(),
	}
	return response, nil
}
