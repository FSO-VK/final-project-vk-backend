package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medbox"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// GetMedicationBox is an interface for getting a Box of medications.
type GetMedicationBox interface {
	Execute(
		ctx context.Context,
		GetMedicationBoxCommand *GetMedicationBoxCommand,
	) (*GetMedicationBoxResponse, error)
}

// GetMedicationBoxService is a service for getting a Box of medications.
type GetMedicationBoxService struct {
	medicationRepo    medication.Repository
	medicationBoxRepo medbox.Repository
	validator         validator.Validator
}

// NewGetMedicationBoxService returns a new GetMedicationBoxService.
func NewGetMedicationBoxService(
	medicationRepo medication.Repository,
	medicationBoxRepo medbox.Repository,
	valid validator.Validator,
) *GetMedicationBoxService {
	return &GetMedicationBoxService{
		medicationRepo:    medicationRepo,
		medicationBoxRepo: medicationBoxRepo,
		validator:         valid,
	}
}

// GetMedicationBoxCommand is a request to get a Box of medications.
type GetMedicationBoxCommand struct {
	UserID string `validate:"required,uuid"`
}

// MedicationBoxItem contains information about one medication in the Box.
type MedicationBoxItem struct {
	ResponseBase
}

// ProducerObject represents object of producer of medication.
type ProducerObject struct {
	Name    string
	Country string
}

// AmountObject is a structure of object of amount of medication.
type AmountObject struct {
	Value float32
	Unit  string
}

// GetMedicationBoxResponse contains a Box of medications.
type GetMedicationBoxResponse struct {
	MedicationBox []*MedicationBoxItem
}

// Execute returns a Box of medications.
func (s *GetMedicationBoxService) Execute(
	ctx context.Context,
	req *GetMedicationBoxCommand,
) (*GetMedicationBoxResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	medBox, err := s.medicationBoxRepo.GetMedicationBox(ctx, userUUID)
	if err != nil {
		if errors.Is(err, medbox.ErrNoMedicationBoxFound) {
			return &GetMedicationBoxResponse{
				MedicationBox: make([]*MedicationBoxItem, 0),
			}, nil
		}
		return nil, fmt.Errorf("failed to get medication box: %w", err)
	}

	items := make([]*MedicationBoxItem, 0, len(medBox.GetMedicationsID()))
	for _, mid := range medBox.GetMedicationsID() {
		med, err := s.medicationRepo.GetByID(ctx, mid)
		if err != nil {
			continue
		}
		items = append(items, &MedicationBoxItem{
			responseBaseMapper(med),
		})
	}

	return &GetMedicationBoxResponse{
		MedicationBox: items,
	}, nil
}
