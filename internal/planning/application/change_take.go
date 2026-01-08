// Package application is a package for application logic of the planning service.
package application

import (
	"context"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/record"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// ChangeTakeMedication is an interface for making a medication taken in any time of the day.
type ChangeTakeMedication interface {
	Execute(
		ctx context.Context,
		cmd *ChangeTakeMedicationCommand,
	) (*ChangeTakeMedicationResponse, error)
}

// ChangeTakeMedicationService is a service for making a medication taken in any time of the day.
type ChangeTakeMedicationService struct {
	recordRepo   record.Repository
	planningRepo plan.Repository
	validator    validator.Validator
}

// NewChangeTakeMedicationService returns a new ChangeTakeMedicationService.
func NewChangeTakeMedicationService(
	recordRepo record.Repository,
	planningRepo plan.Repository,
	valid validator.Validator,
) *ChangeTakeMedicationService {
	return &ChangeTakeMedicationService{
		recordRepo:   recordRepo,
		planningRepo: planningRepo,
		validator:    valid,
	}
}

// ChangeTakeMedicationCommand is a request to change medication take time.
type ChangeTakeMedicationCommand struct {
	PlanID   string `validate:"required,uuid"`
	RecordID string `validate:"required,uuid"`
	UserID   string `validate:"required,uuid"`
	TakenAt  string `validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

// ChangeTakeMedicationResponse is a response to change medication take time.
type ChangeTakeMedicationResponse struct{}

// Execute executes the ChangeTakeMedication command.
func (s *ChangeTakeMedicationService) Execute(
	ctx context.Context,
	req *ChangeTakeMedicationCommand,
) (*ChangeTakeMedicationResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, ErrValidationFail
	}
	parsedPlanID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return nil, ErrValidationFail
	}
	parsedRecordID, err := uuid.Parse(req.RecordID)
	if err != nil {
		return nil, ErrValidationFail
	}

	parsedUser, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, ErrValidationFail
	}

	requestedPlan, err := s.planningRepo.GetByID(ctx, parsedPlanID)
	if err != nil {
		return nil, ErrNoPlan
	}

	parsedTakenAt, err := time.Parse(time.RFC3339, req.TakenAt)
	if err != nil {
		return nil, fmt.Errorf("invalid taking time: %w", err)
	}

	if requestedPlan.UserID() != parsedUser {
		return nil, ErrPlanNotBelongToUser
	}

	requestedRecord, err := s.recordRepo.GetByID(ctx, parsedRecordID)
	if err != nil {
		return nil, ErrNoIntakeRecord
	}

	requestedRecord.MarkTaken(parsedTakenAt)

	err = s.recordRepo.UpdateByID(ctx, requestedRecord)
	if err != nil {
		return nil, err
	}

	return &ChangeTakeMedicationResponse{}, nil
}
