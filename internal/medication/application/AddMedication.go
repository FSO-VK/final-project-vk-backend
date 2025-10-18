// Package application is a package for application logic of the medication service.
package application

import (
	"context"
	"fmt"
	"time"

	medication "github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
)

// AddMedication is an interface for adding a medication.
type AddMedication interface {
	Execute(
		ctx context.Context,
		cmd *AddMedicationCommand,
	) (*AddMedicationResponse, error)
}

// AddMedicationService is a service for adding a medication.
type AddMedicationService struct {
	medicationRepo medication.RepositoryForMedication
	validator      validator.Validator
}

// NewAddMedicationService returns a new AddMedicationService.
func NewAddMedicationService(
	medicationRepo medication.RepositoryForMedication,
	valid validator.Validator,
) *AddMedicationService {
	return &AddMedicationService{
		medicationRepo: medicationRepo,
		validator:      valid,
	}
}

// AddMedicationCommand is a request to add a medication.
type AddMedicationCommand struct {
	Name         string `validate:"required"`
	CategoriesID []uint
	Items        uint   `validate:"required"`
	ItemsUnit    string `validate:"required"`
	Expires      string `validate:"required"`
}

// AddMedicationResponse is a response to add a medication.
type AddMedicationResponse struct {
	ID uint
}

// Execute executes the AddMedication command.
func (s *AddMedicationService) Execute(
	ctx context.Context,
	req *AddMedicationCommand,
) (*AddMedicationResponse, error) {
	err := s.validator.ValidateStruct(req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	expiration, err := time.Parse(time.DateOnly, req.Expires)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiration: %w", err)
	}

	medication := medication.NewMedication(
		req.Name,
		req.Items,
		req.CategoriesID,
		req.ItemsUnit,
		expiration,
	)

	addedMedication, err := s.medicationRepo.Create(ctx, medication)
	if err != nil {
		return nil, fmt.Errorf("failed to add medication: %w", err)
	}

	return &AddMedicationResponse{
		ID: addedMedication.ID,
	}, nil
}
