package credential

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type CredentialType string

const (
	TypeEmail CredentialType = "email"
)

var (
	ErrNotEmailCredentials = errors.New("credentials type is not email")
	ErrEmptyPassword       = errors.New("password is empty")
	ErrSamePassword        = errors.New("password is the same as before")
)

type Credential struct {
	ID             uuid.UUID
	CredentialType CredentialType
	Identifier     string
	Secret         Secret
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewCredential(
	ID uuid.UUID,
	credentialType CredentialType,
	identifier string,
	secret Secret,
	now time.Time,
) *Credential {
	return &Credential{
		ID:             ID,
		CredentialType: credentialType,
		Identifier:     identifier,
		Secret:         secret,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// IsTypeEmail checks if the credential type is email.
func (c *Credential) IsTypeEmail() bool {
	return c.CredentialType == TypeEmail
}

// ChangePassword changes the password of the credential.
func (c *Credential) ChangePassword(secret Secret, now time.Time) error {
	if !c.IsTypeEmail() {
		return ErrNotEmailCredentials
	}

	if secret.GetSecret() == "" {
		return ErrEmptyPassword
	}

	if secret.GetSecret() == c.Secret.GetSecret() {
		return ErrSamePassword
	}
	c.Secret = secret
	c.UpdatedAt = now

	return nil
}
