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

func NewIdentifierEmail(
	plainIdentifier string,
) (*IdentifierEmail, error) {
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

func isValidEmail(email string) bool {
	// by RFC 5322
	pattern := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@` + // local part
		`[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?` + // domain label
		`(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$` // optional subdomains

	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func (s *IdentifierEmail) GetIdentifier() string {
	return s.email
}
