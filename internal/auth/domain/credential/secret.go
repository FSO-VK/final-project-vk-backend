package credential

import (
	"errors"
	"fmt"
	"unicode"
)

// Secret is an interface of Value Object.
type Secret interface {
	GetSecret() string
}

// PasswordHasher is an domain interface.
type PasswordHasher interface {
	Encrypt(password string) (string, error)
	Compare(plainPassword, hashedPassword string) bool
}

const (
	MinPasswordLength = 8
	HasUpper          = true
	HasLower          = true
	HasNumber         = true
)

var (
	ErrPasswordShort = errors.New(
		"password length is less than " + fmt.Sprint(MinPasswordLength),
	)
	ErrPasswordNoUpper  = errors.New("password must contain at least 1 uppercase letter")
	ErrPasswordNoLower  = errors.New("password must contain at least 1 lowercase letter")
	ErrPasswordNoNumber = errors.New("password must contain at least 1 number")
)

// SecretPassword is an implementation of Value Object Secret.
type SecretPassword struct {
	passwordHash string
}

func NewSecretPassword(
	plainPassword string,
	hasher PasswordHasher,
) (*SecretPassword, error) {
	if len(plainPassword) < MinPasswordLength {
		return nil, ErrPasswordShort
	}

	var isUpper, isLower, isNumber bool
	for _, char := range plainPassword {
		if unicode.IsUpper(char) {
			isUpper = true
		} else if unicode.IsLower(char) {
			isLower = true
		} else if unicode.IsNumber(char) {
			isNumber = true
		}
	}

	var err error
	if !isUpper {
		err = errors.Join(err, ErrPasswordNoUpper)
	}
	if !isLower {
		err = errors.Join(err, ErrPasswordNoLower)
	}
	if !isNumber {
		err = errors.Join(err, ErrPasswordNoNumber)
	}
	if err != nil {
		return nil, err
	}

	hashedPassword, err := hasher.Encrypt(plainPassword)
	if err != nil {
		return nil, err
	}
	return &SecretPassword{
		passwordHash: hashedPassword,
	}, nil
}

func (s *SecretPassword) GetSecret() string {
	return s.passwordHash
}
