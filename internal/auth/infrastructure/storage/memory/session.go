package memory

import (
	"context"

	"github.com/FSO-VK/final-project-vk-backend/internal/auth/session/domain"
	"github.com/google/uuid"
)

type SessionStorage struct {
	data Cache[*domain.Session]
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		data: *NewCache[*domain.Session](),
	}
}

func (s *SessionStorage) Create(ctx context.Context, session *domain.Session) error {
	s.data.Set(session.ID.String(), session)
	return nil
}

func (s *SessionStorage) GetByID(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error) {
	session, ok := s.data.Get(sessionID.String())
	if !ok {
		return nil, domain.ErrNoSessionFound
	}
	return session, nil
}

func (s *SessionStorage) Update(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	_, ok := s.data.Get(session.ID.String())
	if !ok {
		return nil, domain.ErrNoSessionFound
	}

	s.data.Set(session.ID.String(), session)
	return session, nil
}

func (s *SessionStorage) Delete(ctx context.Context, sessionID uuid.UUID) error {
	s.data.Delete(sessionID.String())
	return nil
}

