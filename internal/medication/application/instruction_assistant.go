package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application/llm"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application/medreference"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/domain/medication"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

//
//nolint:lll
const RestrictionAnswer = `Извините мы не можем ответить на ваш вопрос основываясь на инструкции выбранного препарата, попробуйте задать другой вопрос или уточните этот.`

// InstructionAssistant is an interface for asking llm about instructions.
type InstructionAssistant interface {
	Execute(
		ctx context.Context,
		cmd *InstructionAssistantCommand,
	) (*InstructionAssistantResponse, error)
}

// InstructionAssistantService is a service for asking llm about instructions.
type InstructionAssistantService struct {
	medicationRepo  medication.Repository
	instructionBot  llm.InstructionBot
	instructionRepo medreference.MedicationReferenceProvider
	validator       validator.Validator
}

// NewInstructionAssistantService returns a new InstructionAssistantService.
func NewInstructionAssistantService(
	medicationRepo medication.Repository,
	instructionBot llm.InstructionBot,
	instructionRepo medreference.MedicationReferenceProvider,
	valid validator.Validator,
) *InstructionAssistantService {
	return &InstructionAssistantService{
		medicationRepo:  medicationRepo,
		instructionBot:  instructionBot,
		instructionRepo: instructionRepo,
		validator:       valid,
	}
}

// InstructionAssistantCommand is a request to ask llm about instructions.
type InstructionAssistantCommand struct {
	UserQuestion string `validate:"required,max=4000"`
	MedicationID string `validate:"required,uuid"`
	UserID       string `validate:"required,uuid"`
}

// InstructionAssistantResponse is a response of llm about instructions.
type InstructionAssistantResponse struct {
	LLMAnswer string
}

// Execute executes the InstructionAssistant command.
func (s *InstructionAssistantService) Execute(
	ctx context.Context,
	req *InstructionAssistantCommand,
) (*InstructionAssistantResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFail, valErr)
	}

	uuidMedicationID, err := uuid.Parse(req.MedicationID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFail, err)
	}

	interestedMedication, err := s.medicationRepo.GetByID(ctx, uuidMedicationID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNoMedication, err)
	}
	barcode := interestedMedication.GetBarCode()

	productInfo, err := s.instructionRepo.GetProductInfo(ctx, barcode)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNoInstruction, err)
	}

	answer, err := s.instructionBot.AskInstructionTwoStep(productInfo.Instruction, req.UserQuestion)
	if err != nil {
		if errors.Is(err, llm.ErrInstructionRestricted) {
			return &InstructionAssistantResponse{
				LLMAnswer: RestrictionAnswer,
			}, nil
		}
		return nil, fmt.Errorf("failed to ask llm: %w", err)
	}

	return &InstructionAssistantResponse{
		LLMAnswer: answer,
	}, nil
}
