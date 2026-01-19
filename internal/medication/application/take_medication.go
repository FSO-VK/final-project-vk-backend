package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medbox"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

var (
	// ErrTakeInvalidUUID represents an error when the uuid is invalid.
	ErrTakeInvalidUUID = errors.New("invalid uuid")
	// ErrTakeInvalidEntity represents an error when the entity is invalid.
	ErrTakeInvalidEntity = errors.New("invalid entity")
	// ErrNotEnoughMedication represents an error when the medication is not enough.
	ErrNotEnoughMedication = errors.New("not enough medication")
)

// TakeMedication is an interface for updating a medication.
type TakeMedication interface {
	Execute(
		ctx context.Context,
		TakeMedicationCommand *TakeMedicationCommand,
	) (*TakeMedicationResponse, error)
}

// TakeMedicationService is a service for updating a medication.
type TakeMedicationService struct {
	medicationRepo    medication.Repository
	medicationBoxRepo medbox.Repository
	validator         validator.Validator
}

// NewTakeMedicationService returns a new TakeMedicationService.
func NewTakeMedicationService(
	medicationRepo medication.Repository,
	medicationBoxRepo medbox.Repository,
	valid validator.Validator,
) *TakeMedicationService {
	return &TakeMedicationService{
		medicationRepo:    medicationRepo,
		medicationBoxRepo: medicationBoxRepo,
		validator:         valid,
	}
}

// TakeMedicationCommand is a request to take a medication (reduce the amount).
type TakeMedicationCommand struct {
	UserID string  `validate:"required,uuid"`
	ID     string  `validate:"required,uuid"`
	Value  float32 `validate:"required,gte=0"`
}

// TakeMedicationResponse is a response to take a medication (reduce the amount).
type TakeMedicationResponse struct {
	// embedded struct
	ResponseBase
}

// Execute takes a medication.
func (s *TakeMedicationService) Execute(
	ctx context.Context,
	req *TakeMedicationCommand,
) (*TakeMedicationResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFail, valErr)
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrTakeInvalidUUID, err)
	}

	uuidUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	oldMedication, err := s.medicationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNoMedication, err)
	}
	if oldMedication.GetAmount().GetValue() < req.Value {
		return nil, fmt.Errorf("%w: %w", ErrNotEnoughMedication, err)
	}

	amount, err := medication.NewMedicationAmount(
		oldMedication.GetAmount().GetValue()-req.Value,
		oldMedication.GetAmount().GetUnit().String(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get amount: %w", err)
	}
	oldMedication.SetAmount(amount)

	savedMedication, err := s.medicationRepo.Update(ctx, oldMedication)
	if err != nil {
		return nil, fmt.Errorf("failed to take medication: %w", err)
	}

	medicationBox, err := s.medicationBoxRepo.GetMedicationBox(ctx, uuidUserID)
	if err != nil {
		return nil, fmt.Errorf("user has no medication box: %w", err)
	}
	err = s.medicationBoxRepo.SetMedicationBox(ctx, medicationBox)
	if err != nil {
		return nil, fmt.Errorf("failed to add medication to box: %w", err)
	}

	return &TakeMedicationResponse{
		responseBaseMapper(savedMedication),
	}, nil
}
