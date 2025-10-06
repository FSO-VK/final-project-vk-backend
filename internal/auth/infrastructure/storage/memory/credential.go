package memory

import (
	"context"

	"github.com/FSO-VK/final-project-vk-backend/internal/auth/domain/credential"
	"github.com/google/uuid"
)

type CredentialStorage struct {
	data *Cache[*credential.Credential]
}

func NewCredentialStorage() *CredentialStorage {
	return &CredentialStorage{
		data: NewCache[*credential.Credential](),
	}
}

func (s *CredentialStorage) Create(ctx context.Context, credential *credential.Credential) error {
	s.data.Set(credential.ID.String(), credential)
	return nil
}

func (s *CredentialStorage) FindByID(
	ctx context.Context,
	credentialID uuid.UUID,
) (*credential.Credential, error) {
	cred, ok := s.data.Get(credentialID.String())
	if !ok {
		return nil, credential.ErrNoCredentialFound
	}
	return cred, nil
}

func (s *CredentialStorage) FindByEmail(
	ctx context.Context,
	email string,
) (*credential.Credential, error) {
	for _, credential := range s.data.data {
		if credential.Identifier == email {
			return credential, nil
		}
	}
	return nil, credential.ErrNoCredentialFound
}
