package credential

import "github.com/google/uuid"

type CredentialType string

const (
	TypeEmail CredentialType = "email"
)

type Credential struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Type       CredentialType
	Identifier string
	Secret     string
}
