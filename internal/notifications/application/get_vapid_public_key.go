// Package application is a package for application logic of the medication service.
package application

import (
	"context"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
)

type PublicKey string

// GetVapidPublicKey is an interface for adding a medication.
type GetVapidPublicKey interface {
	Execute(
		ctx context.Context,
		cmd *GetVapidPublicKeyCommand,
	) (*GetVapidPublicKeyResponse, error)
}

// GetVapidPublicKeyService is a service for getting public key.
type GetVapidPublicKeyService struct {
	publicKey PublicKey
	validator validator.Validator
}

// NewGetVapidPublicKeyService returns a new GetVapidPublicKeyService.
func NewGetVapidPublicKeyService(
	publicKey PublicKey,
	valid validator.Validator,
) *GetVapidPublicKeyService {
	return &GetVapidPublicKeyService{
		publicKey: publicKey,
		validator: valid,
	}
}

// GetVapidPublicKeyCommand is a request to to get public key.
type GetVapidPublicKeyCommand struct {
}

// GetVapidPublicKeyResponse is a response to get public key.
type GetVapidPublicKeyResponse struct {
	PublicKey string
}

// Execute executes the GetVapidPublicKey command.
func (s *GetVapidPublicKeyService) Execute(
	ctx context.Context,
	req *GetVapidPublicKeyCommand,
) (*GetVapidPublicKeyResponse, error) {
	valErr := s.validator.ValidateStruct(req)
	if valErr != nil {
		return nil, fmt.Errorf("failed to validate request: %w", valErr)
	}

	return &GetVapidPublicKeyResponse{
		PublicKey: string(s.publicKey),
	}, nil
}
