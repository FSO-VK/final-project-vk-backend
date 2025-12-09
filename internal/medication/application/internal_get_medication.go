package application

import (
	"context"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medbox"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// InternalGetMedicationByID provides a way to get medication with given medication id.
type InternalGetMedicationByID interface {
	Execute(
		ctx context.Context,
		command *InternalGetMedicationByIDCommand,
	) (*InternalGetMedicationByIDResponse, error)
}

// InternalGetMedicationByIDService is a application service implementing InternalGetMedicationByID interface.
type InternalGetMedicationByIDService struct {
	medicationRepo    medication.Repository
	medicationBoxRepo medbox.Repository
	validator         validator.Validator
}

// NewInternalGetMedicationByIDService creates InternalGetMedicationByIDService.
func NewInternalGetMedicationByIDService(
	medicationRepo medication.Repository,
	medicationBoxRepo medbox.Repository,
	valid validator.Validator,
) *InternalGetMedicationByIDService {
	return &InternalGetMedicationByIDService{
		medicationRepo:    medicationRepo,
		medicationBoxRepo: medicationBoxRepo,
		validator:         valid,
	}
}

// InternalGetMedicationByIDCommand is a command for InternalGetMedicationByID usecase.
type InternalGetMedicationByIDCommand struct {
	ID string `validate:"required,uuid"`
}

// InternalGetMedicationByIDResponse is a response for InternalGetMedicationByID usecase.
type InternalGetMedicationByIDResponse struct {
	ResponseBase
}

// Execute runs InternalGetMedicationByID usecase.
func (s *InternalGetMedicationByIDService) Execute(
	ctx context.Context,
	req *InternalGetMedicationByIDCommand,
) (*InternalGetMedicationByIDResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, ErrValidationFail
	}

	medicationID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, ErrValidationFail
	}

	medication, err := s.medicationRepo.GetByID(ctx, medicationID)
	if err != nil {
		return nil, ErrFailedToGetMedication
	}

	return &InternalGetMedicationByIDResponse{
		responseBaseMapper(medication),
	}, nil
}
