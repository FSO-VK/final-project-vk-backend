// Package application is a package for application logic of the planning service.
package application

import (
	"context"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/plan"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/domain/record"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// CancelMedicationTake is an interface for getting a notification.
type CancelMedicationTake interface {
	Execute(
		ctx context.Context,
		cmd *CancelMedicationTakeCommand,
	) (*CancelMedicationTakeResponse, error)
}

// CancelMedicationTakeService is a service for making a medication taken.
type CancelMedicationTakeService struct {
	recordRepo   record.Repository
	planningRepo plan.Repository
	validator    validator.Validator
}

// NewCancelMedicationTakeService returns a new CancelMedicationTakeService.
func NewCancelMedicationTakeService(
	recordRepo record.Repository,
	planningRepo plan.Repository,
	valid validator.Validator,
) *CancelMedicationTakeService {
	return &CancelMedicationTakeService{
		recordRepo:   recordRepo,
		planningRepo: planningRepo,
		validator:    valid,
	}
}

// CancelMedicationTakeCommand is a request to get a plan.
type CancelMedicationTakeCommand struct {
	PlanID   string `validate:"required,uuid"`
	RecordID string `validate:"required,uuid"`
	UserID   string `validate:"required,uuid"`
}

// CancelMedicationTakeResponse is a response to get a plan.
type CancelMedicationTakeResponse struct{}

// Execute executes the CancelMedicationTake command.
func (s *CancelMedicationTakeService) Execute(
	ctx context.Context,
	req *CancelMedicationTakeCommand,
) (*CancelMedicationTakeResponse, error) {
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

	requestedRecord.MarDefault()

	err = s.recordRepo.UpdateByID(ctx, requestedRecord)
	if err != nil {
		return nil, err
	}

	return &CancelMedicationTakeResponse{}, nil
}
