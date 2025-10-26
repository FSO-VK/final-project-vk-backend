package application

import (
	"context"
	"fmt"

	medication "github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// DeleteMedication is an interface for deleting a medication.
type DeleteMedication interface {
	Execute(
		ctx context.Context,
		DeleteMedicationCommand *DeleteMedicationCommand,
	) (*DeleteMedicationResponse, error)
}

// DeleteMedicationService is a service for deleting a medication.
type DeleteMedicationService struct {
	medicationRepo medication.RepositoryForMedication
	validator      validator.Validator
}

// NewDeleteMedicationService returns a new DeleteMedicationService.
func NewDeleteMedicationService(
	medicationRepo medication.RepositoryForMedication,
	valid validator.Validator,
) *DeleteMedicationService {
	return &DeleteMedicationService{
		medicationRepo: medicationRepo,
		validator:      valid,
	}
}

// DeleteMedicationCommand is a request to delete a medication.
type DeleteMedicationCommand struct {
	ID string `validate:"required"`
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
		return nil, fmt.Errorf("invalid UUID format: %w", err)
	}

	err = s.medicationRepo.Delete(ctx, parsedUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete medication: %w", err)
	}

	return &DeleteMedicationResponse{}, nil
}
