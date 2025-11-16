package main

import (
	"fmt"
	"time"
	"github.com/FSO-VK/final-project-vk-backend/pkg/gigachat"
)

func main() {
	cfg := gigachat.ClientConfig{
		LLMBaseURL:   "https://gigachat.devices.sberbank.ru/api/v1/chat/completions",
		AuthBaseURL:  "https://ngw.devices.sberbank.ru:9443/api/v2/oauth",
		ClientID:     "019a3ae4-0864-79c8-baad-3554efd66556",
		ClientSecret: "c98fbd6b-a132-48fb-ba1e-2a58e3f05369",
		ModelName:    "GigaChat",
		Role:         "user",
		Temperature:  0.1,
		MaxTokens:    2000,
		Timeout:      30 * time.Second,
	}
	provider := gigachat.NewGigachatLLMProvider(cfg)

	// // Demo 1: распознавание планирования из голосового
	// userText := "Я буду принимать Нурофен по 2 таблетки в 7 вечера, в течении месяца"
	// fmt.Println("=== Demo 1: RecognizePlanningFromText ===")
	// fmt.Println("user plan text:", userText)
	// pretty, err := gigachat.RecognizePlanningFromText(provider, userText)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// } else {
	// 	fmt.Println("LLM response:\n", pretty)
	// }

	// Demo 2: two-step question based on instruction
	instruction := gigachat.InstructionDocument
	userQuestion := "Какие противопоказания по принятию и употреблению?"
	fmt.Println("\n=== Demo 2: AskInstructionTwoStep ===")
	fmt.Println("User question:", userQuestion)
	ans, err := gigachat.AskInstructionTwoStep(provider, instruction, userQuestion)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Final LLM answer:\n", ans)
	}
}
