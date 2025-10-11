package credential

import (
	"errors"
	"regexp"
)

// Identifier is an interface of Value Object.
type Identifier interface {
	GetIdentifier() string
}

var (
	ErrInvalidEmail = errors.New("invalid email format")
	ErrEmailTooLong = errors.New("email address too long")
	ErrEmailEmpty   = errors.New("email cannot be empty")
)

// IdentifierEmail is an implementation of Value Object Identifier.
type IdentifierEmail struct {
	email string
}

//nolint:cyclop
func NewIdentifier(
	IdentifierType Credential,
	plainIdentifier string,
) (*IdentifierEmail, error) {
	if !IdentifierType.IsTypeEmail() {
		return nil, ErrNotEmailCredentials
	}
	if plainIdentifier == "" {
		return nil, ErrEmailEmpty
	}

	if len(plainIdentifier) > 254 {
		return nil, ErrEmailTooLong
	}

	if !isValidEmail(plainIdentifier) {
		return nil, ErrInvalidEmail
	}

	return &IdentifierEmail{
		email: plainIdentifier,
	}, nil
}

// by RFC 5322
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`

	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func (s *IdentifierEmail) GetIdentifier() string {
	return s.email
}
