// Package llm is a package for LLM use cases.
package llm

import "errors"

var (
	// ErrLLMInternalFailure is an error when LLM fails.
	ErrLLMInternalFailure = errors.New("llm internal error")
	// ErrInstructionRestricted is an error when question is not about instruction.
	ErrInstructionRestricted = errors.New("instruction restriction")
)

// InstructionBot is an interface for getting instruction advice.
type InstructionBot interface {
	AskInstructionTwoStep(instruction any, userQuestion string) (string, error)
}
