// Package instruction is an application service interface for instruction.
package instruction

import "context"

// Instruction is a data structure for instruction.
type Instruction struct{}

// Service is an interface for instruction application service.
type Service interface {
	GetInstruction(ctx context.Context, barCode string) (*Instruction, error)
}
