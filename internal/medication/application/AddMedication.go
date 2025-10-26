// Package application is a package for application logic of the medication service.
package application

import (
	"context"
	"fmt"
	"time"

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
	medicationRepo medication.RepositoryForMedication
	validator      validator.Validator
}

// NewAddMedicationService returns a new AddMedicationService.
func NewAddMedicationService(
	medicationRepo medication.RepositoryForMedication,
	valid validator.Validator,
) *AddMedicationService {
	return &AddMedicationService{
		medicationRepo: medicationRepo,
		validator:      valid,
	}
}

// AddMedicationCommand is a request to add a medication.
type AddMedicationCommand struct {
	Name              string `validate:"required"`
	InternationalName string
	AmountValue       float32 `validate:"required,gte=0"`
	AmountUnit        string  `validate:"required"`
	ReleaseForm       string  `validate:"required"`

	Group string

	ManufacturerName    string
	ManufacturerCountry string

	ActiveSubstanceName string
	ActiveSubstanceDose float32
	ActiveSubstanceUnit string

	Expires string `validate:"required"`
	Release string

	Commentary string
}

// AddMedicationResponse is a response to add a medication.
type AddMedicationResponse struct {
	AddMedicationCommand

	ID string
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
		release, err := time.Parse(time.DateOnly, req.Release)
		if err != nil {
			return nil, fmt.Errorf("failed to parse release: %w", err)
		}
		expiration = release
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid v7: %w", err)
	}

	drug, err := medication.NewMedicationParse(
		medication.MedicationDraft{
			ID:             id,
			Name:           req.Name,
			ReleaseForm:    req.ReleaseForm,
			AmountValue:    req.AmountValue,
			AmountUnit:     req.AmountUnit,
			ExpirationDate: expiration,

			InternationalName: req.InternationalName,
			Group:             req.Group,
			Manufacturer: medication.MedicationManufacturerDraft{
				Name:    req.ManufacturerName,
				Country: req.ManufacturerCountry,
			},
			ActiveSubstanceName:      req.ActiveSubstanceName,
			ActiveSubstanceDoseValue: req.ActiveSubstanceDose,
			ActiveSubstanceDoseUnit:  req.ActiveSubstanceUnit,
			Commentary:               req.Commentary,
			ReleaseDate:              release,
			CreatedAt:                time.Now(),
			UpdatedAt:                time.Now(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create medication: %w", err)
	}

	addedMedication, err := s.medicationRepo.Create(ctx, drug)
	if err != nil {
		return nil, fmt.Errorf("failed to save medication: %w", err)
	}

	return &AddMedicationResponse{
		ID: addedMedication.ID.String(),
		AddMedicationCommand: AddMedicationCommand{
			Name:                addedMedication.GetName().GetName(),
			InternationalName:   addedMedication.GetInternationalName().GetInternationalName(),
			AmountValue:         addedMedication.GetAmount().GetValue(),
			AmountUnit:          addedMedication.GetAmount().GetUnit().String(),
			ReleaseForm:         addedMedication.GetReleaseForm().String(),
			Group:               addedMedication.GetGroup().GetGroup(),
			ManufacturerName:    addedMedication.GetManufacturer().GetName(),
			ManufacturerCountry: addedMedication.GetManufacturer().GetCountry(),
			ActiveSubstanceName: addedMedication.GetActiveSubstance().GetName(),
			ActiveSubstanceDose: addedMedication.GetActiveSubstance().GetDose().GetValue(),
			ActiveSubstanceUnit: addedMedication.GetActiveSubstance().GetDose().GetUnit().String(),
			Expires:             addedMedication.GetExpirationDate().Format(time.DateOnly),
			Release:             addedMedication.GetReleaseDate().Format(time.DateOnly),
			Commentary:          addedMedication.GetCommentary().GetCommentary(),
		},
	}, nil
}
