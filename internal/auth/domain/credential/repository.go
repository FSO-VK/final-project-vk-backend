package credential

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrNoCredentialFound = errors.New("no credential found")

type CredentialRepository interface {
	Create(ctx context.Context, credential *Credential) error
	FindByID(ctx context.Context, credentialID uuid.UUID) (*Credential, error)
	FindByEmail(ctx context.Context, email string) (*Credential, error)
}
