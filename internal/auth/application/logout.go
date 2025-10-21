package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/auth/domain/session"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

var (
	ErrLogoutValidationFail = errors.New("logout command validation failed")
)

type Logout interface {
	Execute(ctx context.Context, cmd *LogoutCommand) (*LogoutResult, error)
}

type LogoutCommand struct {
	SessionID string `validate:"required,uuid"`
}

type LogoutResult struct {
	SessionID string
}

type LogoutService struct {
	sessionRepo session.SessionRepository
	validator   validator.Validator
}

func NewLogoutService(
	sessionRepo session.SessionRepository,
	validator validator.Validator,
) *LogoutService {
	return &LogoutService{
		sessionRepo: sessionRepo,
		validator:   validator,
	}
}

func (s *LogoutService) Execute(
	ctx context.Context,
	cmd *LogoutCommand,
) (*LogoutResult, error) {
	err := s.validator.ValidateStruct(cmd)
	if err != nil {
		errors.Join(err, ErrLogoutValidationFail)
		return nil, fmt.Errorf("invalid logout command: %w", err)
	}

	sessionIDuuid, err := uuid.Parse(cmd.SessionID)
	if err != nil {
		return nil, fmt.Errorf("parse session id to uuid: %w", err)
	}

	err = s.sessionRepo.Delete(ctx, sessionIDuuid)
	if err != nil {
		return nil, fmt.Errorf("failed to delete session by id: %w", err)
	}

	return &LogoutResult{
		SessionID: sessionIDuuid.String(),
	}, nil
}
