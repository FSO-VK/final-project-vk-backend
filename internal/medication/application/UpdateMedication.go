package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

var (
	ErrUpdateInvalidUUID   = errors.New("invalid uuid")
	ErrUpdateInvalidEntity = errors.New("invalid entity")
)

// UpdateMedication is an interface for updating a medication.
type UpdateMedication interface {
	Execute(
		ctx context.Context,
		UpdateMedicationCommand *UpdateMedicationCommand,
	) (*UpdateMedicationResponse, error)
}

// UpdateMedicationService is a service for updating a medication.
type UpdateMedicationService struct {
	medicationRepo medication.RepositoryForMedication
	validator      validator.Validator
}

// NewUpdateMedicationService returns a new UpdateMedicationService.
func NewUpdateMedicationService(
	medicationRepo medication.RepositoryForMedication,
	valid validator.Validator,
) *UpdateMedicationService {
	return &UpdateMedicationService{
		medicationRepo: medicationRepo,
		validator:      valid,
	}
}

// UpdateMedicationCommand is a request to update a medication.
type UpdateMedicationCommand struct {
	// fields embedded
	AddMedicationCommand

	ID string `validate:"required,uuid"`
}

// UpdateMedicationResponse is a response to update a medication.
type UpdateMedicationResponse struct {
	AddMedicationCommand

	ID string
}

// Execute updates a medication.
func (s *UpdateMedicationService) Execute(
	ctx context.Context,
	req *UpdateMedicationCommand,
) (*UpdateMedicationResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdateInvalidUUID, err)
	}

	oldMedication, err := s.medicationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get medication: %w", err)
	}

	updatedMedication, err := s.updateMedicationEntity(req, oldMedication)
	if err != nil {
		return nil, fmt.Errorf("failed to update medication entity: %w", err)
	}

	savedMedication, err := s.medicationRepo.Update(ctx, updatedMedication)
	if err != nil {
		return nil, fmt.Errorf("failed to update medication: %w", err)
	}

	return &UpdateMedicationResponse{
		ID: savedMedication.ID.String(),
		AddMedicationCommand: AddMedicationCommand{
			Name:                savedMedication.GetName().GetName(),
			InternationalName:   savedMedication.GetInternationalName().GetInternationalName(),
			AmountValue:         savedMedication.GetAmount().GetValue(),
			AmountUnit:          savedMedication.GetAmount().GetUnit().String(),
			ReleaseForm:         savedMedication.GetReleaseForm().String(),
			Group:               savedMedication.GetGroup().GetGroup(),
			ManufacturerName:    savedMedication.GetManufacturer().GetName(),
			ManufacturerCountry: savedMedication.GetManufacturer().GetCountry(),
			ActiveSubstanceName: savedMedication.GetActiveSubstance().GetName(),
			ActiveSubstanceDose: savedMedication.GetActiveSubstance().GetDose().GetValue(),
			ActiveSubstanceUnit: savedMedication.GetActiveSubstance().GetDose().GetUnit().String(),
			Expires:             savedMedication.GetExpirationDate().Format(time.DateOnly),
			Release:             savedMedication.GetReleaseDate().Format(time.DateOnly),
			Commentary:          savedMedication.GetCommentary().GetCommentary(),
		},
	}, nil
}

func (u *UpdateMedicationService) updateMedicationEntity(
	req *UpdateMedicationCommand,
	oldMedication *medication.Medication,
) (*medication.Medication, error) {
	var allErrors error

	name, err := medication.NewMedicationName(req.Name)
	allErrors = errors.Join(allErrors, err)

	internationalName, err := medication.NewMedicationInternationalName(
		req.InternationalName)
	allErrors = errors.Join(allErrors, err)

	amount, err := medication.NewMedicationAmount(
		req.AmountValue,
		req.AmountUnit,
	)
	allErrors = errors.Join(allErrors, err)

	group, err := medication.NewMedicationGroup(req.Group)
	allErrors = errors.Join(allErrors, err)

	manufacturer, err := medication.NewMedicationManufacturer(
		req.ManufacturerName,
		req.ManufacturerCountry,
	)
	allErrors = errors.Join(allErrors, err)

	activeSubstance, err := medication.NewMedicationActiveSubstance(
		req.ActiveSubstanceName,
		req.ActiveSubstanceDose,
		req.ActiveSubstanceUnit,
	)
	allErrors = errors.Join(allErrors, err)

	expiration, err := time.Parse(time.DateOnly, req.Expires)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiration: %w", err)
	}

	release, err := time.Parse(time.DateOnly, req.Release)
	if err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	if allErrors != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdateInvalidEntity, allErrors)
	}

	oldMedication.SetName(name)
	oldMedication.SetInternationalName(internationalName)
	oldMedication.SetAmount(amount)
	oldMedication.SetGroup(group)
	oldMedication.SetManufacturer(manufacturer)
	oldMedication.SetActiveSubstance(activeSubstance)

	oldMedication.SetUpdatedAt(time.Now())
	oldMedication.SetReleaseDate(release)
	oldMedication.SetExpirationDate(expiration)

	return oldMedication, nil
}
