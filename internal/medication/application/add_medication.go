// Package application is a package for application logic of the medication service.
package application

import (
	"context"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medbox"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
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
	medicationRepo    medication.Repository
	medicationBoxRepo medbox.Repository
	validator         validator.Validator
}

// NewAddMedicationService returns a new AddMedicationService.
func NewAddMedicationService(
	medicationRepo medication.Repository,
	medicationBoxRepo medbox.Repository,
	valid validator.Validator,
) *AddMedicationService {
	return &AddMedicationService{
		medicationRepo:    medicationRepo,
		medicationBoxRepo: medicationBoxRepo,
		validator:         valid,
	}
}

// AddMedicationCommand is a request to add a medication.
type AddMedicationCommand struct {
	// embedded struct
	CommandBase

	UserID string `validate:"required,uuid"`
}

// AddMedicationResponse is a response to add a medication.
type AddMedicationResponse struct {
	// embedded struct
	ResponseBase
}

// Execute executes the AddMedication command.
func (s *AddMedicationService) Execute(
	ctx context.Context,
	req *AddMedicationCommand,
) (*AddMedicationResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}

	expiration, err := time.Parse(time.DateOnly, req.Expires)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiration: %w", err)
	}

	var release time.Time
	if req.Release != "" {
		release, err = time.Parse(time.DateOnly, req.Release)
		if err != nil {
			return nil, fmt.Errorf("failed to parse release: %w", err)
		}
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid v7: %w", err)
	}

	drug, err := medication.Parse(
		medication.MedicationDraft{
			ID:             id,
			Name:           req.Name,
			ReleaseForm:    req.ReleaseForm,
			AmountValue:    req.AmountValue,
			AmountUnit:     req.AmountUnit,
			ExpirationDate: expiration,

			InternationalName: req.InternationalName,
			Group:             req.Group,
			Manufacturer: medication.ManufacturerDraft{
				Name:    req.ManufacturerName,
				Country: req.ManufacturerCountry,
			},
			ActiveSubstance: MapActiveSubstanceToDraft(req.ActiveSubstance),
			Commentary:      req.Commentary,
			ReleaseDate:     release,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create medication: %w", err)
	}

	addedMedication, err := repositoryModifications(ctx, s, req, drug)
	if err != nil {
		return nil, fmt.Errorf("failed to add medication: %w", err)
	}

	return &AddMedicationResponse{
		responseBaseMapper(addedMedication),
	}, nil
}

func repositoryModifications(
	ctx context.Context, s *AddMedicationService,
	req *AddMedicationCommand, drug *medication.Medication,
) (*medication.Medication, error) {
	addedMedication, err := s.medicationRepo.Create(ctx, drug)
	if err != nil || addedMedication == nil {
		return nil, fmt.Errorf("failed to save medication: %w", err)
	}
	uuidUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	medicationBox, err := s.medicationBoxRepo.GetMedicationBox(ctx, uuidUserID)
	if err != nil {
		medicationBoxDraft := medbox.NewMedicationBox(uuidUserID)
		medicationBox, err = s.medicationBoxRepo.CreateMedicationBox(ctx, medicationBoxDraft)
		if err != nil {
			return nil, fmt.Errorf("failed to create medication box: %w", err)
		}
	}
	medicationBox.AddMedication(addedMedication.GetID())
	err = s.medicationBoxRepo.SetMedicationBox(ctx, medicationBox)
	if err != nil {
		return nil, fmt.Errorf("failed to add medication to box: %w", err)
	}

	return addedMedication, nil
}
