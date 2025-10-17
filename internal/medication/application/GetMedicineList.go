package application

import (
	"context"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medicine"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
)

// GetMedicineList is an interface for getting a list of medicines.
type GetMedicineList interface {
	Execute(
		ctx context.Context,
		GetMedicineListCommand *GetMedicineListCommand,
	) (*GetMedicineListResponse, error)
}

// GetMedicineListService is a service for getting a list of medicines.
type GetMedicineListService struct {
	medicineRepo medicine.RepositoryForMedication
	validator    validator.Validator
}

// NewGetMedicineListService returns a new GetMedicineListService.
func NewGetMedicineListService(
	medicineRepo medicine.RepositoryForMedication,
	valid validator.Validator,
) *GetMedicineListService {
	return &GetMedicineListService{
		medicineRepo: medicineRepo,
		validator:    valid,
	}
}

// GetMedicineListCommand is a request to get a list of medicines.
type GetMedicineListCommand struct{}

// MedicineListItem contains information about one medicine in the list.
type MedicineListItem struct {
	ID        uint
	Name      string
	Items     uint
	ItemsUnit string
	Expires   string
}

// GetMedicineListResponse contains a list of medicines.
type GetMedicineListResponse struct {
	List []*MedicineListItem
}

// Execute returns a list of medicines.
func (s *GetMedicineListService) Execute(
	ctx context.Context,
	_ *GetMedicineListCommand,
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
