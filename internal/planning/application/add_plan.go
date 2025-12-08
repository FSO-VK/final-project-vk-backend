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
	"github.com/teambition/rrule-go"
)

// ErrUnsupportedRrule is an error when rrule is unsupported.
var ErrUnsupportedRrule = errors.New("rrule is unsupported")

// AddPlan is an interface for adding a notification.
type AddPlan interface {
	Execute(
		ctx context.Context,
		cmd *AddPlanCommand,
	) (*AddPlanResponse, error)
}

// AddPlanService is a service for creating a subscription.
type AddPlanService struct {
	planningRepo plan.Repository
	validator    validator.Validator
}

// NewAddPlanService returns a new AddPlanService.
func NewAddPlanService(
	planningRepo plan.Repository,
	valid validator.Validator,
) *AddPlanService {
	return &AddPlanService{
		planningRepo: planningRepo,
		validator:    valid,
	}
}

// AddPlanCommand is a request to add a plan.
type AddPlanCommand struct {
	MedicationID   string   `validate:"required,uuid"`
	UserID         string   `validate:"required,uuid"`
	AmountValue    float64  `validate:"required,gte=0"`
	AmountUnit     string   `validate:"required"`
	Condition      string   `validate:"omitempty,max=300"`
	StartDate      string   `validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate        string   `validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	RecurrenceRule []string `validate:"required"`
}

// AddPlanResponse is a response to add a plan.
type AddPlanResponse struct {
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

// Execute executes the AddPlan command.
func (s *AddPlanService) Execute(
	ctx context.Context,
	req *AddPlanCommand,
) (*AddPlanResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, ErrValidationFail
	}
	parsedUser, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, ErrValidationFail
	}
	parsedMedicationID, err := uuid.Parse(req.MedicationID)
	if err != nil {
		return nil, ErrValidationFail
	}
	newPlan, err := createPlan(req, parsedUser, parsedMedicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}
	err = s.planningRepo.Save(ctx, newPlan)
	if err != nil {
		return nil, fmt.Errorf("failed to save plan: %w", err)
	}
	amountValue, amountUnit := newPlan.Dosage()

	response := &AddPlanResponse{
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

func createPlan(req *AddPlanCommand,
	userID uuid.UUID,
	medicationID uuid.UUID,
) (*plan.Plan, error) {
	dosage, err := plan.NewDosage(
		req.AmountValue,
		req.AmountUnit,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid dosage: %w", err)
	}

	parsedStart, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid course start: %w", err)
	}
	parsedEnd, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid course end: %w", err)
	}
	if len(req.RecurrenceRule) == 0 {
		return nil, ErrUnsupportedRrule
	}
	rules := make([]*rrule.RRule, 0, len(req.RecurrenceRule))

	for _, ruleStr := range req.RecurrenceRule {
		rule, err := rrule.StrToRRule(ruleStr)
		if err != nil {
			return nil, ErrUnsupportedRrule
		}
		rules = append(rules, rule)
	}
	schedule, err := plan.NewSchedule(parsedStart, parsedEnd, rules)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule: %w", err)
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate uuid: %w", err)
	}
	newPlan, err := plan.NewPlan(
		id,
		medicationID,
		userID,
		dosage,
		schedule,
		req.Condition,
		time.Now(),
		time.Now(),
	)
	return newPlan, err
}
