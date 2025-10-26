package application

import (
	"context"
	"fmt"
	"time"

	medication "github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
)

// GetMedicationList is an interface for getting a list of medications.
type GetMedicationList interface {
	Execute(
		ctx context.Context,
		GetMedicationListCommand *GetMedicationListCommand,
	) (*GetMedicationListResponse, error)
}

// GetMedicationListService is a service for getting a list of medications.
type GetMedicationListService struct {
	medicationRepo medication.RepositoryForMedication
	validator      validator.Validator
}

// NewGetMedicationListService returns a new GetMedicationListService.
func NewGetMedicationListService(
	medicationRepo medication.RepositoryForMedication,
	valid validator.Validator,
) *GetMedicationListService {
	return &GetMedicationListService{
		medicationRepo: medicationRepo,
		validator:      valid,
	}
}

// GetMedicationListCommand is a request to get a list of medications.
type GetMedicationListCommand struct{}

// MedicationListItem contains information about one medication in the list.
type MedicationListItem struct {
	ID        uint
	Name      string
	Items     uint
	ItemsUnit string
	Expires   string
}

// GetMedicationListResponse contains a list of medications.
type GetMedicationListResponse struct {
	List []*MedicationListItem
}

// Execute returns a list of medications.
func (s *GetMedicationListService) Execute(
	ctx context.Context,
	req *GetMedicationListCommand,
) (*GetMedicationListResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}
	medications, err := s.medicationRepo.GetListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get medication list: %w", err)
	}

	listItems := make([]*MedicationListItem, 0)
	for _, medication := range medications {
		listItems = append(listItems, &MedicationListItem{
			ID:        medication.ID,
			Name:      medication.Name,
			Items:     medication.Items,
			ItemsUnit: medication.ItemsUnit,
			Expires:   medication.ExpiresAt.Format(time.DateOnly),
		})
	}

	return &GetMedicationListResponse{
		List: listItems,
	}, nil
}
