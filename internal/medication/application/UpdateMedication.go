package application

import (
	"context"
	"fmt"
	"time"

	medication "github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
)

// UpdateMedication is an interface for updating a medication.
type UpdateMedication interface {
	Execute(
		ctx context.Context,
		UpdateMedicationCommand *UpdateMedicationCommand,
	) (*UpdateMedicationResponse, error)
}

// UpdateMedicationService is a service for updating a medication.
type UpdateMedicationService struct {
	medicationRepo medication.RepositoryForMedication
	validator      validator.Validator
}

// NewUpdateMedicationService returns a new UpdateMedicationService.
func NewUpdateMedicationService(
	medicationRepo medication.RepositoryForMedication,
	valid validator.Validator,
) *UpdateMedicationService {
	return &UpdateMedicationService{
		medicationRepo: medicationRepo,
		validator:      valid,
	}
}

// UpdateMedicationCommand is a request to update a medication.
type UpdateMedicationCommand struct {
	ID           uint
	Name         string `validate:"required"`
	CategoriesID []uint
	Items        uint   `validate:"required"`
	ItemsUnit    string `validate:"required"`
	Expires      string `validate:"required"`
}

// UpdateMedicationResponse is a response to update a medication.
type UpdateMedicationResponse struct {
	ID           uint
	Name         string
	CategoriesID []uint
	Items        uint
	ItemsUnit    string
	Expires      string
}

// Execute updates a medication.
func (s *UpdateMedicationService) Execute(
	ctx context.Context,
	req *UpdateMedicationCommand,
) (*UpdateMedicationResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}

	expiration, err := time.Parse(time.DateOnly, req.Expires)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiration: %w", err)
	}

	id := req.ID
	oldMedication, err := s.medicationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get medication: %w", err)
	}

	medication := medication.NewMedication(
		req.Name,
		req.Items,
		req.CategoriesID,
		req.ItemsUnit,
		expiration,
	)
	medication.ID = oldMedication.ID

	updatedMedication, err := s.medicationRepo.Update(ctx, medication)
	if err != nil {
		return nil, fmt.Errorf("failed to update medication: %w", err)
	}

	return &UpdateMedicationResponse{
		ID:           updatedMedication.ID,
		Name:         updatedMedication.Name,
		CategoriesID: updatedMedication.CategoriesID,
		Items:        updatedMedication.Items,
		ItemsUnit:    updatedMedication.ItemsUnit,
		Expires:      updatedMedication.Expires.Format(time.DateOnly),
	}, nil
}
