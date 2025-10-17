// Package application is a package for application logic of the medication service.
package application

import (
	"context"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medicine"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
)

// AddMedicine is an interface for adding a medicine.
type AddMedicine interface {
	Execute(
		ctx context.Context,
		cmd *AddMedicineCommand,
	) (*AddMedicineResponse, error)
}

// AddMedicineService is a service for adding a medicine.
type AddMedicineService struct {
	medicineRepo medicine.RepositoryForMedication
	validator    validator.Validator
}

// NewAddMedicineService returns a new AddMedicineService.
func NewAddMedicineService(
	medicineRepo medicine.RepositoryForMedication,
	valid validator.Validator,
) *AddMedicineService {
	return &AddMedicineService{
		medicineRepo: medicineRepo,
		validator:    valid,
	}
}

// AddMedicineCommand is a request to add a medicine.
type AddMedicineCommand struct {
	Name         string `validate:"required"`
	CategoriesID []uint
	Items        uint   `validate:"required"`
	ItemsUnit    string `validate:"required"`
	Expires      string `validate:"required"`
}

// AddMedicineResponse is a response to add a medicine.
type AddMedicineResponse struct {
	ID uint
}

// Execute executes the AddMedicine command.
func (s *AddMedicineService) Execute(
	ctx context.Context,
	req *AddMedicineCommand,
) (*AddMedicineResponse, error) {
	// err := s.validator.ValidateStruct(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to validate request: %w", err)
	// }

	expiration, err := time.Parse(time.DateOnly, req.Expires)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiration: %w", err)
	}

	medicine := medicine.NewMedicine(
		req.Name,
		req.Items,
		req.CategoriesID,
		req.ItemsUnit,
		expiration,
	)

	addedMedicine, err := s.medicineRepo.Create(ctx, medicine)
	if err != nil {
		return nil, fmt.Errorf("failed to add medicine: %w", err)
	}

	return &AddMedicineResponse{
		ID: addedMedicine.ID,
	}, nil
}
