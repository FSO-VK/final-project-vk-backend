package credential

import "github.com/google/uuid"

type CredentialType string

const (
	TypeEmail CredentialType = "email"
)

type Credential struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	CredentialType CredentialType
	Identifier     string
	Secret         Secret
}

func NewCredential(
	ID uuid.UUID,
	userID uuid.UUID,
	credentialType CredentialType,
	identifier string,
	secret Secret,
) *Credential {
	return &Credential{
		ID:             ID,
		UserID:         userID,
		CredentialType: credentialType,
		Identifier:     identifier,
		Secret:         secret,
	}
}

func (c *Credential) IsTypeEmail() bool {
	return c.CredentialType == TypeEmail
}
