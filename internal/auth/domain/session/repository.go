package session

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNoSessionFound = errors.New("no session found")
)

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, sessionID uuid.UUID) (*Session, error)
	Update(ctx context.Context, session *Session) (*Session, error)
	Delete(ctx context.Context, sessionID uuid.UUID) error
}
