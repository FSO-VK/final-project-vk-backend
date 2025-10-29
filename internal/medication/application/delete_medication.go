package application

import (
	"context"
	"errors"
	"fmt"

	medication "github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medbox"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// ErrDeleteInvalidUUIDFormat represents an error when the uuid is invalid.
var ErrDeleteInvalidUUIDFormat = errors.New("invalid UUID format")

// DeleteMedication is an interface for deleting a medication.
type DeleteMedication interface {
	Execute(
		ctx context.Context,
		DeleteMedicationCommand *DeleteMedicationCommand,
	) (*DeleteMedicationResponse, error)
}

// DeleteMedicationService is a service for deleting a medication.
type DeleteMedicationService struct {
	medicationRepo    medication.RepositoryForMedication
	medicationBoxRepo medbox.RepositoryForMedicationBox
	validator         validator.Validator
}

// NewDeleteMedicationService returns a new DeleteMedicationService.
func NewDeleteMedicationService(
	medicationRepo medication.RepositoryForMedication,
	medicationBoxRepo medbox.RepositoryForMedicationBox,
	valid validator.Validator,
) *DeleteMedicationService {
	return &DeleteMedicationService{
		medicationRepo:    medicationRepo,
		medicationBoxRepo: medicationBoxRepo,
		validator:         valid,
	}
}

// DeleteMedicationCommand is a request to delete a medication.
type DeleteMedicationCommand struct {
	UserID string `validate:"required"`
	ID     string `validate:"required"`
}

// DeleteMedicationResponse is a response to delete a medication.
type DeleteMedicationResponse struct{}

// Execute deletes a medication.
func (s *DeleteMedicationService) Execute(
	ctx context.Context,
	req *DeleteMedicationCommand,
) (*DeleteMedicationResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}

	parsedUUID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDeleteInvalidUUIDFormat, err)
	}

	err = s.medicationRepo.Delete(ctx, parsedUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete medication: %w", err)
	}

	uuidUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	medicationBox, err := s.medicationBoxRepo.GetMedicationBox(ctx, uuidUserID)
	if err != nil {
		return nil, fmt.Errorf("user does not have a medication box: %w", err)
	}
	medicationBox.MedicationsID = removeFromSlice(medicationBox.MedicationsID, parsedUUID)
	err = s.medicationBoxRepo.SetMedicationBox(ctx, medicationBox)
	if err != nil {
		return nil, fmt.Errorf("failed to add medication to box: %w", err)
	}

	return &DeleteMedicationResponse{}, nil
}

func removeFromSlice(slice []uuid.UUID, id uuid.UUID) []uuid.UUID {
	for i, item := range slice {
		if item == id {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
