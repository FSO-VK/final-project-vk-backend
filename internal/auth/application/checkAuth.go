package application

import (
	"context"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/auth/domain/session"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

type CheckAuth interface {
	Execute(ctx context.Context, checkAuthCommand *CheckAuthCommand) (*CheckAuthResult, error)
}

// CheckAuthCommand represents the command to check authentication.
type CheckAuthCommand struct {
	SessionID string
}

// CheckAuthResult represents the result of a check authentication operation.
type CheckAuthResult struct {
	SessionID       string
	IsAuthenticated bool
	ExpiresAt       time.Time
}

type CheckAuthService struct {
	sessionRepo session.SessionRepository
	validator   validator.Validator
}

func NewCheckAuthService(
	sessionRepo session.SessionRepository,
	valid validator.Validator,
) *CheckAuthService {
	return &CheckAuthService{
		sessionRepo: sessionRepo,
		validator:   valid,
	}
}

func (s *CheckAuthService) Execute(
	ctx context.Context,
	checkAuthCommand *CheckAuthCommand,
) (*CheckAuthResult, error) {
	err := s.validator.ValidateStruct(checkAuthCommand)
	if err != nil {
		return nil, fmt.Errorf("invalid check auth command: %w", err)
	}

	sessionId, err := uuid.Parse(checkAuthCommand.SessionID)
	if err != nil {
		return nil, fmt.Errorf("parse session id to uuid: %w", err)
	}

	session, err := s.sessionRepo.GetByID(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by id: %w", err)
	}

	if session.IsExpired() || session.IsRevoked() {
		return &CheckAuthResult{
			SessionID:       session.ID.String(),
			IsAuthenticated: false,
			ExpiresAt:       session.ExpiresAt,
		}, nil
	}

	err = session.Refresh()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}

	_, err = s.sessionRepo.Update(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return &CheckAuthResult{
		SessionID:       session.ID.String(),
		IsAuthenticated: true,
		ExpiresAt:       session.ExpiresAt,
	}, nil
}
