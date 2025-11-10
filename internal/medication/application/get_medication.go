package application

import (
	"context"
	"errors"
	"slices"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medbox"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// GetMedicationByID provides a way to get medication with given medication id.
type GetMedicationByID interface {
	Execute(ctx context.Context, command *GetMedicationByIDCommand) (*GetMedicationByIDResponse, error)
}

// GetMedicationByIDService is a application service implementing GetMedicationById interface.
type GetMedicationByIDService struct {
	medicationRepo    medication.Repository
	medicationBoxRepo medbox.Repository
	validator         validator.Validator
}

// NewGetMedicationByIDService creates GetMedicationByIdService.
func NewGetMedicationByIDService(
	medicationRepo medication.Repository,
	medicationBoxRepo medbox.Repository,
	valid validator.Validator,
) *GetMedicationByIDService {
	return &GetMedicationByIDService{
		medicationRepo:    medicationRepo,
		medicationBoxRepo: medicationBoxRepo,
		validator:         valid,
	}
}

// GetMedicationByIDCommand is a command for GetMedicationByID usecase.
type GetMedicationByIDCommand struct {
	UserID string `validate:"required,uuid"`
	ID     string `validate:"required,uuid"`
}

// GetMedicationByIDResponse is a response for GetMedicationByID usecase.
type GetMedicationByIDResponse struct {
	ResponseBase
}

var (
	// ErrFailedToGetMedication occurs when repository fails to get medication.
	ErrFailedToGetMedication = errors.New("failed to get medication")
)

// Execute runs GetMedicationbyID usecase.
func (s *GetMedicationByIDService) Execute(
	ctx context.Context,
	req *GetMedicationByIDCommand,
) (*GetMedicationByIDResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, ErrValidationFail
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, ErrValidationFail
	}

	medicationID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, ErrValidationFail
	}

	medBox, err := s.medicationBoxRepo.GetMedicationBox(ctx, userID)
	if err != nil {
		if errors.Is(err, medbox.ErrNoMedicationBoxFound) {
			// Since userID is basically MedboxID, this is not possible.
			// But we must process this error
			return nil, ErrNoMedication
		}
	}
	medicationsInTheBox := medBox.GetMedicationsID()
	if contains := slices.Contains(medicationsInTheBox, medicationID); !contains {
		return nil, ErrNoMedication
	}

	medication, err := s.medicationRepo.GetByID(ctx, medicationID)
	if err != nil {
		return nil, ErrFailedToGetMedication
	}

	return &GetMedicationByIDResponse{
		responseBaseMapper(medication),
	}, nil
}
