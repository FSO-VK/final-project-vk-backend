package application

import (
	"context"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medicine"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
)

// UpdateMedicine is an interface for updating a medicine.
type UpdateMedicine interface {
	Execute(
		ctx context.Context,
		UpdateMedicineCommand *UpdateMedicineCommand,
	) (*UpdateMedicineResponse, error)
}

// UpdateMedicineService is a service for updating a medicine.
type UpdateMedicineService struct {
	medicineRepo medicine.RepositoryForMedication
	validator    validator.Validator
}

// NewUpdateMedicineService returns a new UpdateMedicineService.
func NewUpdateMedicineService(
	medicineRepo medicine.RepositoryForMedication,
	valid validator.Validator,
) *UpdateMedicineService {
	return &UpdateMedicineService{
		medicineRepo: medicineRepo,
		validator:    valid,
	}
}

// UpdateMedicineCommand is a request to update a medicine.
type UpdateMedicineCommand struct {
	ID           uint
	Name         string `validate:"required"`
	CategoriesID []uint
	Items        uint   `validate:"required"`
	ItemsUnit    string `validate:"required"`
	Expires      string `validate:"required"`
}

// UpdateMedicineResponse is a response to update a medicine.
type UpdateMedicineResponse struct {
	ID           uint
	Name         string
	CategoriesID []uint
	Items        uint
	ItemsUnit    string
	Expires      string
}

// Execute updates a medicine.
func (s *UpdateMedicineService) Execute(
	ctx context.Context,
	req *UpdateMedicineCommand,
) (*UpdateMedicineResponse, error) {
	// err := s.validator.ValidateStruct(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to validate request: %w", err)
	// }

	expiration, err := time.Parse(time.DateOnly, req.Expires)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiration: %w", err)
	}

	id := req.ID
	oldMedicine, err := s.medicineRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get medicine: %w", err)
	}

	medicine := medicine.NewMedicine(
		req.Name,
		req.Items,
		req.CategoriesID,
		req.ItemsUnit,
		expiration,
	)
	medicine.ID = oldMedicine.ID

	updatedMedicine, err := s.medicineRepo.Update(ctx, medicine)
	if err != nil {
		return nil, fmt.Errorf("failed to update medicine: %w", err)
	}

	return &UpdateMedicineResponse{
		ID:           updatedMedicine.ID,
		Name:         updatedMedicine.Name,
		CategoriesID: updatedMedicine.CategoriesID,
		Items:        updatedMedicine.Items,
		ItemsUnit:    updatedMedicine.ItemsUnit,
		Expires:      updatedMedicine.Expires.Format(time.DateOnly),
	}, nil
}
