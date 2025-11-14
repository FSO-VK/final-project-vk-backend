// Package llmchatbot is a package for LLM chat bot.
package llmchatbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"text/template"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application/medreference"
	"github.com/FSO-VK/final-project-vk-backend/pkg/llm"
)

var (
	// ErrEmptyResponse is returned when the response body is empty or contains no data.
	ErrEmptyResponse = errors.New("empty response")
	// ErrWithSystemPrompt is returned when the access token is missing in the response.
	ErrWithSystemPrompt = errors.New("failed to get token")
	// ErrInvalidInstruction is returned when the instruction is invalid.
	ErrInvalidInstruction = errors.New("invalid instruction")
)

// LLMChatBot is a service for getting a Box of medications.
type LLMChatBot struct {
	llmProvider llm.Provider
}

// NewLLMChatBot returns a new LLMChatBot.
func NewLLMChatBot(
	llmProvider llm.Provider,
) *LLMChatBot {
	return &LLMChatBot{
		llmProvider: llmProvider,
	}
}

// SelectInstructionFieldPrompt is a prompt for the LLM to select the instruction field.
type SelectInstructionFieldPrompt struct {
	UserQuestion      string
	InstructionFields string
}

// InstructionConsultationPrompt is a prompt for the LLM to get the instruction consultation.
type InstructionConsultationPrompt struct {
	UserQuestion    string
	InstructionText string
}

// AskInstructionTwoStep asks the LLM a question about the instruction.
func (s *LLMChatBot) AskInstructionTwoStep(
	instruction medreference.Instruction,
	userQuestion string,
) (string, error) {
	template, err := template.ParseFiles(
		"./internal/medication/infrastructure/llm_chat_bot/templates/select_instruction_field.tmpl",
	)
	if err != nil {
		return "", ErrWithSystemPrompt
	}

	instructionFields, err := newEmptyInstructionJSON()
	if err != nil {
		return "", err
	}

	data := SelectInstructionFieldPrompt{
		UserQuestion:      userQuestion,
		InstructionFields: instructionFields,
	}

	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil {
		return "", ErrWithSystemPrompt
	}

	firstPrompt := buf.String()
	resp1, err := s.llmProvider.Query(firstPrompt)
	if err != nil {
		return "", err
	}

	if resp1 == "" {
		return "", ErrEmptyResponse
	}
	instructionPart, err := getFieldValue(instruction, resp1)
	if err != nil {
		return "", err
	}

	templateSecond, err := template.ParseFiles(
		"./internal/medication/infrastructure/llm_chat_bot/templates/instruction_consultation.tmpl",
	)
	if err != nil {
		return "", ErrWithSystemPrompt
	}
	dataSecond := InstructionConsultationPrompt{
		UserQuestion:    userQuestion,
		InstructionText: instructionPart,
	}

	var bufSecond bytes.Buffer
	if err := templateSecond.Execute(&bufSecond, dataSecond); err != nil {
		return "", ErrWithSystemPrompt
	}

	lastPrompt := buf.String()
	finalResponse, err := s.llmProvider.Query(lastPrompt)
	if err != nil {
		return "", err
	}
	return finalResponse, nil
}

func getFieldValue(doc interface{}, fieldName string) (string, error) {
	r := reflect.ValueOf(doc)

	if r.Kind() != reflect.Struct {
		return "", ErrInvalidInstruction
	}

	field := r.FieldByName(fieldName)
	if !field.IsValid() {
		return "", ErrInvalidInstruction
	}

	return field.String(), nil
}

func newEmptyInstructionJSON() (string, error) {
	emptyInstruction := medreference.Instruction{
		Nozologies:             []medreference.Nozology{},
		ClPhPointers:           []medreference.ClPhPointer{},
		PharmInfluence:         "",
		PharmKinetics:          "",
		Dosage:                 "",
		OverDosage:             "",
		Interaction:            "",
		Lactation:              "",
		SideEffects:            "",
		UsingIndication:        "",
		UsingCounterIndication: "",
		SpecialInstruction:     "",
		RenalInsuf:             "",
		HepatoInsuf:            "",
		ElderlyInsuf:           "",
		ChildInsuf:             "",
	}
	//nolint:musttag // nolint because we don't need actual json tags, it is just for giving LLM names of structure fields
	instructionJSON, err := json.Marshal(emptyInstruction)
	if err != nil {
		return "", err
	}
	return string(instructionJSON), nil
}
