package llmchatbot

// InstructionAssistantConfig is configuration for llm assistant.
type InstructionAssistantConfig struct {
	SelectInstructionFieldPromptPath string `koanf:"select_instruction_field_prompt_path"`
	ConsultingPromptPath             string `koanf:"consulting_prompt_path"`
}