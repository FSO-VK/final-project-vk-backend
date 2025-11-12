package gigachat

import (
	"bytes"
	"fmt"
	"html/template"
)

type UserPlan struct {
	UserText string
}

// RecognizePlanningFromText: single-step extraction of med schedule from user text.
func RecognizePlanningFromText(provider *GigachatLLMProvider, userText string) (string, error) {
	template, err := template.ParseFiles("./pkg/gigachat/templates/add_planning_prompt.tmpl")
	if err != nil {
		return "", ErrWithSystemPrompt
	}
	data := UserPlan{UserText: userText}
	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil {
		return "", ErrWithSystemPrompt
	}

	fullPrompt := buf.String()
	resp, err := provider.Query(fullPrompt)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return resp, nil
}
