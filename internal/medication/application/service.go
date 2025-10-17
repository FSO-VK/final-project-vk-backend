// Package application is a package for application logic of the medication service.
package application

import (
	"context"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medicine"
)

// MedicineServiceProvider is a service provider for medicine.
type MedicineServiceProvider struct {
	medicineRepo medicine.RepositoryForMedication
}

// NewMedicineServiceProvider creates a new medicine service provider.
func NewMedicineServiceProvider(
	medicineRepo medicine.RepositoryForMedication,
) *MedicineServiceProvider {
	return &MedicineServiceProvider{
		medicineRepo: medicineRepo,
	}
}

// AddMedicine adds a medicine.
func (s *MedicineServiceProvider) AddMedicine(
	ctx context.Context,
	req *AddMedicineRequest,
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

// UpdateMedicine updates a medicine.
func (s *MedicineServiceProvider) UpdateMedicine(
	ctx context.Context,
	req *UpdateMedicineRequest,
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

// DeleteMedicine deletes a medicine.
func (s *MedicineServiceProvider) DeleteMedicine(
	ctx context.Context,
	req *DeleteMedicineRequest,
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

// GetMedicineList returns a list of medicines.
func (s *MedicineServiceProvider) GetMedicineList(
	ctx context.Context,
	_ *GetMedicineListRequest,
) (*GetMedicineListResponse, error) {
	medicines, err := s.medicineRepo.GetListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get medicine list: %w", err)
	}

	listItems := make([]*MedicineListItem, 0)
	for _, medicine := range medicines {
		listItems = append(listItems, &MedicineListItem{
			ID:        medicine.ID,
			Name:      medicine.Name,
			Items:     medicine.Items,
			ItemsUnit: medicine.ItemsUnit,
			Expires:   medicine.Expires.Format(time.DateOnly),
		})
	}

	return &GetMedicineListResponse{
		List: listItems,
	}, nil
}
