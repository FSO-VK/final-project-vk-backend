package session

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type SessionStatus string

const (
	StatusActive  SessionStatus = "active"
	StatusRevoked SessionStatus = "revoked"
)

var (
	ErrSessionRefreshExpired = errors.New("cannot refresh expired session")
	ErrSessionRevoked        = errors.New("session already revoked")
)

var SessionDuration = 24 * time.Hour

type Session struct {
	ID           uuid.UUID
	CredentialID uuid.UUID
	Status       SessionStatus
	LastLoginAt  time.Time
	ExpiresAt    time.Time
}

func NewSession(credentialID uuid.UUID) *Session {
	return &Session{
		ID:           uuid.New(),
		CredentialID: credentialID,
		Status:       StatusActive,
		LastLoginAt:  time.Now(),
		ExpiresAt:    time.Now().Add(SessionDuration),
	}
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s *Session) Refresh() error {
	if s.IsRevoked() {
		return ErrSessionRevoked
	}

	if s.IsExpired() {
		return ErrSessionRefreshExpired
	}

	s.LastLoginAt = time.Now()
	s.ExpiresAt = time.Now().Add(SessionDuration)
	return nil
}

func (s *Session) Revoke() {
	s.Status = StatusRevoked
}

func (s *Session) IsRevoked() bool {
	return s.Status == StatusRevoked
}
