package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/auth/domain/session"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

var (
	ErrNoValidSession = errors.New("user is not authenticated")
	ErrNoSessionFound = errors.New("no session found")
)

type CheckAuth interface {
	Execute(ctx context.Context, checkAuthCommand *CheckAuthCommand) (*CheckAuthResult, error)
}

// CheckAuthCommand represents the command to check authentication.
type CheckAuthCommand struct {
	SessionID string `validate:"required"`
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
		return nil, ErrNoValidSession
	}

	sessionId, err := uuid.Parse(checkAuthCommand.SessionID)
	if err != nil {
		return nil, ErrNoValidSession
	}

	userSession, err := s.sessionRepo.GetByID(ctx, sessionId)
	if err != nil {
		if errors.Is(err, session.ErrNoSessionFound) {
			return nil, ErrNoSessionFound
		}
		return nil, fmt.Errorf("fail with db: %w", err)
	}

	if userSession.IsExpired() || userSession.IsRevoked() {
		return &CheckAuthResult{
			SessionID:       userSession.ID.String(),
			IsAuthenticated: false,
			ExpiresAt:       userSession.ExpiresAt,
		}, nil
	}

	err = userSession.Refresh()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}

	_, err = s.sessionRepo.Update(ctx, userSession)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return &CheckAuthResult{
		SessionID:       userSession.ID.String(),
		IsAuthenticated: true,
		ExpiresAt:       userSession.ExpiresAt,
	}, nil
}
