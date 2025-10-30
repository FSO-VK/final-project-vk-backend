package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medbox"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

var (
	// ErrUpdateInvalidUUID represents an error when the uuid is invalid.
	ErrUpdateInvalidUUID = errors.New("invalid uuid")
	// ErrUpdateInvalidEntity represents an error when the entity is invalid.
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
	medicationRepo    medication.Repository
	medicationBoxRepo medbox.Repository
	validator         validator.Validator
}

// NewUpdateMedicationService returns a new UpdateMedicationService.
func NewUpdateMedicationService(
	medicationRepo medication.Repository,
	medicationBoxRepo medbox.Repository,
	valid validator.Validator,
) *UpdateMedicationService {
	return &UpdateMedicationService{
		medicationRepo:    medicationRepo,
		medicationBoxRepo: medicationBoxRepo,
		validator:         valid,
	}
}

// UpdateMedicationCommand is a request to update a medication.
type UpdateMedicationCommand struct {
	// fields embedded
	CommandBase

	UserID string `validate:"required,uuid"`
	ID     string `validate:"required,uuid"`
}

// UpdateMedicationResponse is a response to update a medication.
type UpdateMedicationResponse struct {
	// embedded struct
	ResponseBase
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

	uuidUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	medicationBox, err := s.medicationBoxRepo.GetMedicationBox(ctx, uuidUserID)
	if err != nil {
		return nil, fmt.Errorf("user has no medication box: %w", err)
	}
	medicationBox.AddMedication(savedMedication.GetID())
	err = s.medicationBoxRepo.SetMedicationBox(ctx, medicationBox)
	if err != nil {
		return nil, fmt.Errorf("failed to add medication to box: %w", err)
	}

	return &UpdateMedicationResponse{
		responseBaseMapper(savedMedication),
	}, nil
}

func (s *UpdateMedicationService) updateMedicationEntity(
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
