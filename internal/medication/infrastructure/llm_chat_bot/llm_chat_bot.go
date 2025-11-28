// Package llmchatbot is a package for LLM chat bot.
package llmchatbot

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"text/template"

	"github.com/FSO-VK/final-project-vk-backend/pkg/llm"
)

var (
	// ErrEmptyResponse is returned when the response body is empty or contains no data.
	ErrEmptyResponse = errors.New("empty response")
	// ErrWithSystemPrompt is returned when the system prompt is missing.
	ErrWithSystemPrompt = errors.New("err with system prompt")
	// ErrInvalidInstruction is returned when the instruction is invalid.
	ErrInvalidInstruction = errors.New("invalid instruction")
)

// LLMChatBot is a service for getting instruction advice.
type LLMChatBot struct {
	llmProvider llm.Provider
	conf        InstructionAssistantConfig
}

// NewLLMChatBot returns a new LLMChatBot.
func NewLLMChatBot(
	llmProvider llm.Provider,
	conf InstructionAssistantConfig,
) *LLMChatBot {
	return &LLMChatBot{
		llmProvider: llmProvider,
		conf:        conf,
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
	instruction any,
	userQuestion string,
) (string, error) {
	selectFieldTemplate, err := template.ParseFiles(s.conf.SelectInstructionFieldPromptPath)
	if err != nil {
		return "", ErrWithSystemPrompt
	}

	instructionFields, err := getFieldNamesList(instruction)
	if err != nil {
		return "", err
	}

	data := SelectInstructionFieldPrompt{
		UserQuestion:      userQuestion,
		InstructionFields: instructionFields,
	}

	var buf bytes.Buffer
	if err := selectFieldTemplate.Execute(&buf, data); err != nil {
		return "", ErrWithSystemPrompt
	}

	selectInstructionFieldPrompt := buf.String()
	LLMChosenInstructionField, err := s.llmProvider.Query(selectInstructionFieldPrompt)
	if err != nil {
		return "", err
	}

	instructionPart, err := getFieldValue(instruction, LLMChosenInstructionField)
	if err != nil {
		return "", err
	}

	consultTemplate, err := template.ParseFiles(s.conf.ConsultingPromptPath)
	if err != nil {
		return "", ErrWithSystemPrompt
	}
	instructionConsultationPromptData := InstructionConsultationPrompt{
		UserQuestion:    userQuestion,
		InstructionText: instructionPart,
	}

	var bufSecond bytes.Buffer
	if err := consultTemplate.Execute(&bufSecond, instructionConsultationPromptData); err != nil {
		return "", ErrWithSystemPrompt
	}

	LLMFinalResponse := bufSecond.String()
	finalResponse, err := s.llmProvider.Query(LLMFinalResponse)
	if err != nil {
		return "", err
	}
	return finalResponse, nil
}

func getFieldValue(doc interface{}, fieldName string) (string, error) {
	fieldName = strings.ToUpper(fieldName[:1]) + fieldName[1:]
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

func getFieldNamesList(instr any) (string, error) {
	t := reflect.TypeOf(instr)
	if t.Kind() != reflect.Struct {
		return "", ErrInvalidInstruction
	}

	numFields := t.NumField()
	fieldNames := make([]string, 0, numFields)

	for i := range numFields {
		fieldNames = append(fieldNames, t.Field(i).Name)
	}

	return strings.Join(fieldNames, "\n"), nil
}
