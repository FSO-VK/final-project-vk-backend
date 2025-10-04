package application

import (
	"context"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medicine/medicine"
)

type MedicineServiceProvider struct {
	medicineRepo medicine.MedicineRepository
}

func NewMedicineServiceProvider(medicineRepo medicine.MedicineRepository) *MedicineServiceProvider {
	return &MedicineServiceProvider{
		medicineRepo: medicineRepo,
	}
}

func (s *MedicineServiceProvider) AddMedicine(
	ctx context.Context,
	req *AddMedicineRequest,
) (*AddMedicineResponse, error) {
	// TODO: add validator

	expiration, err := time.Parse(time.RFC3339, req.Expires)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiration: %w", err)
	}

	medicine := medicine.NewMedicine(
		req.Name,
		req.Items,
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

func (s *MedicineServiceProvider) UpdateMedicine(
	ctx context.Context,
	req *UpdateMedicineRequest,
) (*UpdateMedicineResponse, error) {
	// TODO: add validator

	expiration, err := time.Parse(time.RFC3339, req.Expires)
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
		req.ItemsUnit,
		expiration,
	)
	medicine.ID = oldMedicine.ID

	updatedMedicine, err := s.medicineRepo.Update(ctx, medicine)
	if err != nil {
		return nil, fmt.Errorf("failed to update medicine: %w", err)
	}

	return &UpdateMedicineResponse{
		ID: updatedMedicine.ID,
	}, nil
}

func (s *MedicineServiceProvider) DeleteMedicine(
	ctx context.Context,
	req *DeleteMedicineRequest,
) (*DeleteMedicineResponse, error) {
	// TODO: add validator

	err := s.medicineRepo.Delete(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete medicine: %w", err)
	}

	return &DeleteMedicineResponse{}, nil
}

func (s *MedicineServiceProvider) GetMedicineList(
	ctx context.Context,
	req *GetMedicineListRequest,
) (*GetMedicineListResponse, error) {
	// TODO: add validator

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
			Expires:   medicine.Expires.Format(time.RFC3339),
		})
	}

	return &GetMedicineListResponse{
		List: listItems,
	}, nil
}
