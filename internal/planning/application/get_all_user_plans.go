package application

import (
	"context"
	"errors"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// GetAllPlans is an interface for getting a notification.
type GetAllPlans interface {
	Execute(
		ctx context.Context,
		cmd *GetAllPlansCommand,
	) (*GetAllPlansResponse, error)
}

// GetAllPlansService is a service for creating a subscription.
type GetAllPlansService struct {
	planningRepo plan.Repository
	validator    validator.Validator
}

// NewGetAllPlansService returns a new GetAllPlansService.
func NewGetAllPlansService(
	planningRepo plan.Repository,
	valid validator.Validator,
) *GetAllPlansService {
	return &GetAllPlansService{
		planningRepo: planningRepo,
		validator:    valid,
	}
}

// GetAllPlansCommand is a request to get a plan.
type GetAllPlansCommand struct {
	UserID string `validate:"required,uuid"`
}

// PlanItem is a plan item.
type PlanItem struct {
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

// GetAllPlansResponse is a response to get a plan.
type GetAllPlansResponse struct {
	Plans []*PlanItem
}

// Execute executes the GetAllPlans command.
func (s *GetAllPlansService) Execute(
	ctx context.Context,
	req *GetAllPlansCommand,
) (*GetAllPlansResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, ErrValidationFail
	}

	parsedUser, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, ErrValidationFail
	}

	userPlans, err := s.planningRepo.UserPlans(ctx, parsedUser)
	if err != nil && !(errors.Is(err, plan.ErrNoPlanFound)) {
		return nil, ErrNoPlan
	}

	plansList := make([]*PlanItem, 0, len(userPlans))

	for _, onePlan := range userPlans {
		amountValue, amountUnit := onePlan.Dosage()
		plansList = append(plansList, &PlanItem{
			ID:             onePlan.ID().String(),
			MedicationID:   onePlan.MedicationID().String(),
			UserID:         onePlan.UserID().String(),
			AmountValue:    amountValue,
			AmountUnit:     amountUnit,
			Condition:      onePlan.Condition(),
			StartDate:      onePlan.CourseStart().Format(time.RFC3339),
			EndDate:        onePlan.CourseEnd().Format(time.RFC3339),
			RecurrenceRule: onePlan.ScheduleIcal(),
		})
	}

	response := &GetAllPlansResponse{
		Plans: plansList,
	}
	return response, nil
}
