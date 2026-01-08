// Package application is a package for application logic of the planning service.
package application

import (
	"context"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/record"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// TakeMedication is an interface for getting a notification.
type TakeMedication interface {
	Execute(
		ctx context.Context,
		cmd *TakeMedicationCommand,
	) (*TakeMedicationResponse, error)
}

// TakeMedicationService is a service for making a medication taken.
type TakeMedicationService struct {
	recordRepo   record.Repository
	planningRepo plan.Repository
	validator    validator.Validator
}

// NewTakeMedicationService returns a new TakeMedicationService.
func NewTakeMedicationService(
	recordRepo record.Repository,
	planningRepo plan.Repository,
	valid validator.Validator,
) *TakeMedicationService {
	return &TakeMedicationService{
		recordRepo:   recordRepo,
		planningRepo: planningRepo,
		validator:    valid,
	}
}

// TakeMedicationCommand is a request to get a plan.
type TakeMedicationCommand struct {
	PlanID   string `validate:"required,uuid"`
	RecordID string `validate:"required,uuid"`
	UserID   string `validate:"required,uuid"`
}

// TakeMedicationResponse is a response to get a plan.
type TakeMedicationResponse struct{}

// Execute executes the TakeMedication command.
func (s *TakeMedicationService) Execute(
	ctx context.Context,
	req *TakeMedicationCommand,
) (*TakeMedicationResponse, error) {
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

	if requestedPlan.UserID() != parsedUser {
		return nil, ErrPlanNotBelongToUser
	}

	requestedRecord, err := s.recordRepo.GetByID(ctx, parsedRecordID)
	if err != nil {
		return nil, ErrNoIntakeRecord
	}

	requestedRecord.MarkTaken(requestedRecord.PlannedTime())

	err = s.recordRepo.UpdateByID(ctx, requestedRecord)
	if err != nil {
		return nil, err
	}

	return &TakeMedicationResponse{}, nil
}
