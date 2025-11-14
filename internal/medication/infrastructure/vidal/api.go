// Package vidal is an implementation of instruction application service.
package vidal

import (
	"context"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application/instruction"
)

type Service struct {
	storage Storage
	client  Client
}

func NewService(storage Storage, client Client) *Service {
	return &Service{storage: storage, client: client}
}

func (s *Service) GetInstruction(ctx context.Context, barCode string) (*instruction.Instruction, error) {
	return nil, nil
}