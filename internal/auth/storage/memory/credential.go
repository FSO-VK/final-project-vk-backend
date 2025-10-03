package memory

import (
	"context"

	"github.com/FSO-VK/final-project-vk-backend/internal/auth/credential/domain"
	"github.com/google/uuid"
)

type CredentialStorage struct {
	data *Cache[*domain.Credential]
}

func NewCredentialStorage() *CredentialStorage {
	return &CredentialStorage{
		data: NewCache[*domain.Credential](),
	}
}

func (s *CredentialStorage) Create(ctx context.Context, credential *domain.Credential) error {
	s.data.Set(credential.ID.String(), credential)
	return nil
}

func (s *CredentialStorage) FindByID(ctx context.Context, credentialID uuid.UUID) (*domain.Credential, error) {
	cred, ok := s.data.Get(credentialID.String())
	if !ok {
		return nil, domain.ErrNoCredentialFound
	}
	return cred, nil
}

func (s *CredentialStorage) FindByEmail(ctx context.Context, email string) (*domain.Credential, error) {
	for _, credential := range s.data.data {
		if credential.Identifier == email {
			return credential, nil
		}
	}
	return nil, domain.ErrNoCredentialFound
}
