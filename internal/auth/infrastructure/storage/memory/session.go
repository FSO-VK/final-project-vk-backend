package memory

import (
	"context"

	"github.com/FSO-VK/final-project-vk-backend/internal/auth/domain/session"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/cache"
	"github.com/google/uuid"
)

type SessionStorage struct {
	data cache.Cache[*session.Session]
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		data: *cache.NewCache[*session.Session](),
	}
}

func (s *SessionStorage) Create(ctx context.Context, session *session.Session) error {
	s.data.Set(session.ID.String(), session)
	return nil
}

func (s *SessionStorage) GetByID(
	ctx context.Context,
	sessionID uuid.UUID,
) (*session.Session, error) {
	sess, ok := s.data.Get(sessionID.String())
	if !ok {
		return nil, session.ErrNoSessionFound
	}
	return sess, nil
}

func (s *SessionStorage) Update(
	ctx context.Context,
	sess *session.Session,
) (*session.Session, error) {
	_, ok := s.data.Get(sess.ID.String())
	if !ok {
		return nil, session.ErrNoSessionFound
	}

	s.data.Set(sess.ID.String(), sess)
	return sess, nil
}

func (s *SessionStorage) Delete(ctx context.Context, sessionID uuid.UUID) error {
	s.data.Delete(sessionID.String())
	return nil
}
