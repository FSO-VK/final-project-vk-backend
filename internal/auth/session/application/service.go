package application

import (
	"context"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/auth/session/domain"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

type SessionServiceProvider struct {
	valid       validator.Validator
	sessionRepo domain.SessionRepository
}

func (s *SessionServiceProvider) CheckAuth(
	ctx context.Context,
	checkReq *CheckAuthRequest,
) (*CheckAuthResponse, error) {
	err := s.valid.ValidateStruct(checkReq)
	if err != nil {
		return nil, fmt.Errorf("check auth request: %w", err)
	}

	sessionID, err := uuid.Parse(checkReq.SessionID)
	if err != nil {
		return nil, fmt.Errorf("parse session id to uuid: %w", err)
	}

	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by id: %w", err)
	}

	if session.IsExpired() || session.IsRevoked() {
		return &CheckAuthResponse{
			SessionID:    sessionID.String(),
			IsAuthorized: false,
			ExpiresAt:    session.ExpiresAt,
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

	return &CheckAuthResponse{
		SessionID:    sessionID.String(),
		IsAuthorized: true,
		ExpiresAt:    session.ExpiresAt,
	}, nil
}
