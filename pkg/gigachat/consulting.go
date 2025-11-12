package gigachat

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

type SelectInstructionFieldPrompt struct {
	UserQuestion string
}

// AskInstructionTwoStep: two-step Q/A on medication instruction
func AskInstructionTwoStep(provider *GigachatLLMProvider, instruction Document, userQuestion string) (string, error) {
	// STEP 1: ask which excerpt is needed
	template, err := template.ParseFiles("./pkg/gigachat/templates/select_instruction_field.tmpl")
	if err != nil {
		return "", ErrWithSystemPrompt
	}
	data := SelectInstructionFieldPrompt{UserQuestion: userQuestion}

	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil {
		return "", ErrWithSystemPrompt
	}

	firstPrompt := buf.String()
	resp1, err := provider.Query(firstPrompt)
	if err != nil {
		fmt.Println("err", err)
		return "", err
	}
	fmt.Println("LLM need that field of instruction:", resp1)

	if resp1 == "" {
		fmt.Println("No excerpt suggested; querying full instruction")
		return "", ErrEmptyResponse
	}
	// STEP 2: ask the question on the excerpt
	instructionPart, err := getFieldValueByJSONTag(instruction, resp1)
	if err != nil {
		return "", err
	}
	return askFinal(provider, instructionPart, userQuestion)
}

// func getFieldValue(doc interface{}, fieldName string) (string, error) {
// 	r := reflect.ValueOf(doc)

// 	if r.Kind() != reflect.Struct {
// 		return "", fmt.Errorf("ожидалась структура, получен %v", r.Kind())
// 	}

// 	field := r.FieldByName(fieldName)
// 	if !field.IsValid() {
// 		return "", fmt.Errorf("поле %s не найдено", fieldName)
// 	}

//		return field.String(), nil
//	}
func getFieldValueByJSONTag(doc interface{}, jsonFieldName string) (string, error) {
	r := reflect.ValueOf(doc)

	if r.Kind() != reflect.Struct {
		return "", fmt.Errorf("ожидалась структура, получен %v", r.Kind())
	}

	t := r.Type()

	// Ищем поле по JSON тегу
	for i := 0; i < r.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")

		// Убираем опции из JSON тега (например, "omitempty")
		if strings.Contains(jsonTag, ",") {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		if jsonTag == jsonFieldName {
			fieldValue := r.Field(i)
			return fieldValue.String(), nil
		}
	}

	return "", fmt.Errorf("поле с JSON тегом %s не найдено", jsonFieldName)
}

type InstructionConsultationPrompt struct {
	UserQuestion        string
	InstructionText string
}

func askFinal(provider *GigachatLLMProvider, instructionSnippet, question string) (string, error) {

	template, err := template.ParseFiles("./pkg/gigachat/templates/instruction_consultation.tmpl")
	if err != nil {
		return "", ErrWithSystemPrompt
	}
	data := InstructionConsultationPrompt{UserQuestion: question, InstructionText: instructionSnippet}

	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil {
		return "", ErrWithSystemPrompt
	}

	lastPrompt := buf.String()
	finalResponse, err := provider.Query(lastPrompt)
	if err != nil {
		return "", err
	}
	return finalResponse, nil
}
