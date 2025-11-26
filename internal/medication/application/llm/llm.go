package llm

// InstructionBot is an interface for getting instruction advice.
type InstructionBot interface {
	AskInstructionTwoStep(instruction any, userQuestion string) (string, error)
}
