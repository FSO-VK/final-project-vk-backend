package application

import (
	"context"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medicine"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
)

// DeleteMedicine is an interface for deleting a medicine.
type DeleteMedicine interface {
	Execute(
		ctx context.Context,
		DeleteMedicineCommand *DeleteMedicineCommand,
	) (*DeleteMedicineResponse, error)
}

// DeleteMedicineService is a service for deleting a medicine.
type DeleteMedicineService struct {
	medicineRepo medicine.RepositoryForMedication
	validator    validator.Validator
}

// NewDeleteMedicineService returns a new DeleteMedicineService.
func NewDeleteMedicineService(
	medicineRepo medicine.RepositoryForMedication,
	valid validator.Validator,
) *DeleteMedicineService {
	return &DeleteMedicineService{
		medicineRepo: medicineRepo,
		validator:    valid,
	}
}

// DeleteMedicineCommand is a request to delete a medicine.
type DeleteMedicineCommand struct {
	ID uint
}

// DeleteMedicineResponse is a response to delete a medicine.
type DeleteMedicineResponse struct{}

// Execute deletes a medicine.
func (s *DeleteMedicineService) Execute(
	ctx context.Context,
	req *DeleteMedicineCommand,
) (*DeleteMedicineResponse, error) {
	// err := s.validator.ValidateStruct(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to validate request: %w", err)
	// }

	err := s.medicineRepo.Delete(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete medicine: %w", err)
	}

	return &DeleteMedicineResponse{}, nil
}
