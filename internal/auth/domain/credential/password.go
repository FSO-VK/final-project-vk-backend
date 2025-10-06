package credential

import (
	"errors"
	"fmt"
	"unicode"
)

const (
	MinPasswordLength = 8
	HasUpper          = true
	HasLower          = true
	HasNumber         = true
)

var (
	ErrPasswordShort    = errors.New("password length is less than " + fmt.Sprint(MinPasswordLength))
	ErrPasswordNoUpper  = errors.New("password must contain at least 1 uppercase letter")
	ErrPasswordNoLower  = errors.New("password must contain at least 1 lowercase letter")
	ErrPasswordNoNumber = errors.New("password must contain at least 1 number")
)

type Password struct {
	value string
}

func NewPassword(value string) (*Password, error) {
	if len(value) < MinPasswordLength {
		return nil, ErrPasswordShort
	}

	var isUpper, isLower, isNumber bool
	for _, char := range value {
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

	return &Password{
		value: value,
	}, nil
}

func (p *Password) String() string {
	return p.value
}

type HashedPassword struct {
	value string
}

func NewHashedPassword(value string) *HashedPassword {
	return &HashedPassword{
		value: value,
	}
}

func (hp *HashedPassword) String() string {
	return hp.value
}
