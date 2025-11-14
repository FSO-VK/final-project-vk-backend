// Package vidal is an implementation of instruction application service.
package vidal

import (
	"context"
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application/instruction"
	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

type Service struct {
	storage Storage
	client  Client
}

func NewService(storage Storage, client Client) *Service {
	return &Service{storage: storage, client: client}
}

func (s *Service) GetInstruction(ctx context.Context, barCode string) (*instruction.Instruction, error) {
	err := validation.EAN13(barCode)
	if err != nil {
		return nil, instruction.ErrBadBarCode
	}

	model, err := s.storage.GetProduct(ctx, barCode)
	if err != nil {
		if !errors.Is(err, ErrNoProduct) {
			return nil, fmt.Errorf("get product from storage: %w", err)
		}
	} else {
		return s.productInfoToInstruction(model)
	}

	clientResponse, err := s.client.GetInstruction(ctx, barCode)
	if err != nil {
		return nil, fmt.Errorf("external service: %w", err)
	}

	model, err = s.clientResponseToModel(clientResponse)
	if err != nil {
		return nil, fmt.Errorf("convert client response to model: %w", err)
	}

	err = s.storage.SaveProduct(ctx, model)
	if err != nil {
		return nil, fmt.Errorf("save product to s: %w", err)
	}

	return s.productInfoToInstruction(model)
}

func (s *Service) productInfoToInstruction(product *StorageModel) (*instruction.Instruction, error) {
	instruction := &instruction.Instruction{}
	return instruction, nil
}

func (s *Service) clientResponseToModel(clientResponse *ClientResponse) (*StorageModel, error) {
	barCodes := make([]string, len(clientResponse.ProductPackages))
	for i, pack := range clientResponse.ProductPackages {
		barCodes[i] = pack.BarCode
	}

	model := &StorageModel{
		Product:  clientResponse.Product,
		BarCodes: barCodes,
	}
	return model, nil
}
