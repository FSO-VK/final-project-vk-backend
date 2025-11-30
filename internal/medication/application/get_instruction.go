package application

import (
	"context"
	"errors"
	"slices"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application/medreference"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medbox"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

// GetInstructionByMedicationID provides a way to get instruction with given instruction id.
type GetInstructionByMedicationID interface {
	Execute(
		ctx context.Context,
		command *GetInstructionByMedicationIDCommand,
	) (*GetInstructionByMedicationIDResponse, error)
}

// GetInstructionByMedicationIDService is a application service implementing GetInstructionByMedicationID interface.
type GetInstructionByMedicationIDService struct {
	medicationRepo    medication.Repository
	medicationBoxRepo medbox.Repository
	instructionRepo   medreference.MedicationReferenceProvider
	validator         validator.Validator
}

// NewGetInstructionByMedicationIDService creates GetInstructionByMedicationIDService.
func NewGetInstructionByMedicationIDService(
	medicationRepo medication.Repository,
	medicationBoxRepo medbox.Repository,
	instructionRepo medreference.MedicationReferenceProvider,
	valid validator.Validator,
) *GetInstructionByMedicationIDService {
	return &GetInstructionByMedicationIDService{
		medicationRepo:    medicationRepo,
		medicationBoxRepo: medicationBoxRepo,
		instructionRepo:   instructionRepo,
		validator:         valid,
	}
}

// GetInstructionByMedicationIDCommand is a command for GetInstructionByMedicationID usecase.
type GetInstructionByMedicationIDCommand struct {
	UserID string `validate:"required,uuid"`
	ID     string `validate:"required,uuid"`
}

// Nosology is an illness.
type Nosology struct {
	Code string
	Name string
}

// ClPhPointer is a clinical-pharmacological pointer.
type ClPhPointer struct {
	Code string
	Name string
}

// GetInstructionByMedicationIDResponse is a response for GetInstructionByMedicationID usecase.
type GetInstructionByMedicationIDResponse struct {
	Nosologies             []Nosology
	ClPhPointers           []ClPhPointer
	PharmInfluence         string
	PharmKinetics          string
	Dosage                 string
	OverDosage             string
	Interaction            string
	Lactation              string
	SideEffects            string
	UsingIndication        string
	UsingCounterIndication string
	SpecialInstruction     string
	RenalInsuf             string
	HepatoInsuf            string
	ElderlyInsuf           string
	ChildInsuf             string
}

// ErrFailedToGetMedication occurs when repository fails to get instruction.
var ErrFailedToGetInstruction = errors.New("failed to get instruction")

// Execute runs GetInstructionByMedicationID usecase.
func (s *GetInstructionByMedicationIDService) Execute(
	ctx context.Context,
	req *GetInstructionByMedicationIDCommand,
) (*GetInstructionByMedicationIDResponse, error) {
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

	productInfo, err := s.instructionRepo.GetProductInfo(ctx, medication.GetBarCode())
	if err != nil {
		return nil, ErrFailedToGetInstruction
	}

	return &GetInstructionByMedicationIDResponse{
		Nosologies:             convertToNosology(productInfo.Instruction.Nozologies),
		ClPhPointers:           convertToClPhPointers(productInfo.Instruction.ClPhPointers),
		PharmInfluence:         productInfo.Instruction.PharmInfluence,
		PharmKinetics:          productInfo.Instruction.PharmKinetics,
		Dosage:                 productInfo.Instruction.Dosage,
		OverDosage:             productInfo.Instruction.OverDosage,
		Interaction:            productInfo.Instruction.Interaction,
		Lactation:              productInfo.Instruction.Lactation,
		SideEffects:            productInfo.Instruction.SideEffects,
		UsingIndication:        productInfo.Instruction.UsingIndication,
		UsingCounterIndication: productInfo.Instruction.UsingCounterIndication,
		SpecialInstruction:     productInfo.Instruction.SpecialInstruction,
		RenalInsuf:             productInfo.Instruction.RenalInsuf,
		HepatoInsuf:            productInfo.Instruction.HepatoInsuf,
		ElderlyInsuf:           productInfo.Instruction.ElderlyInsuf,
		ChildInsuf:             productInfo.Instruction.ChildInsuf,
	}, nil
}

func convertToNosology(substances []medreference.Nozology) []Nosology {
	result := make([]Nosology, len(substances))
	for i, v := range substances {
		result[i] = Nosology{
			Code: v.Code,
			Name: v.Name,
		}
	}
	return result
}

func convertToClPhPointers(substances []medreference.ClPhPointer) []ClPhPointer {
	result := make([]ClPhPointer, len(substances))
	for i, v := range substances {
		result[i] = ClPhPointer{
			Code: v.Code,
			Name: v.Name,
		}
	}
	return result
}
