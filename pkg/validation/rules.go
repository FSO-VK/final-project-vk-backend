// Package validation contains validation rules.
package validation

import (
	"errors"
	"fmt"
	"unicode"
)

var (
	ErrValueRequired = errors.New("value can't be empty")
	ErrValueShort    = errors.New("value is too short")
	ErrValueLong     = errors.New("value is too long")
	ErrValueNegative = errors.New("value can't be negative")
	ErrValueFormat   = errors.New("value has invalid format")
)

func Required(value string) error {
	if value == "" {
		return ErrValueRequired
	}
	return nil
}

func MinLength(value string, length int) error {
	if len(value) < length {
		return fmt.Errorf("%w: can't be less than %d", ErrValueShort, length)
	}
	return nil
}

func MaxLength(value string, length int) error {
	if len(value) > length {
		return fmt.Errorf("%w, can't be longer than %d", ErrValueLong, length)
	}
	return nil
}

type Numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 |
		uint16 | uint32 | uint64 | float32 | float64
}

func Positive[T Numeric](value T) error {
	if value < 0 {
		return ErrValueNegative
	}
	return nil
}

// ----- Specific GS1-related validators ----- //
// RequiredGTIN ensures GTIN is exactly 14 digits (numeric).
func RequiredGTIN(value string) error {
	if err := MinLength(value, 14); err != nil {
		return err
	}
	if err := MaxLength(value, 14); err != nil {
		return err
	}
	for _, r := range value {
		if !unicode.IsDigit(r) {
			return fmt.Errorf("%w: gtin must contain only digits", ErrValueFormat)
		}
	}
	return nil
}

// RequiredSerial ensures SerialNumber is exactly 13 characters (letters/digits allowed).
// Accepts Latin letters (upper/lower) and digits only.
func RequiredSerial(value string) error {
	if err := MinLength(value, 13); err != nil {
		return err
	}
	if err := MaxLength(value, 13); err != nil {
		return err
	}
	for _, r := range value {
		if !unicode.IsDigit(r) && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return fmt.Errorf("%w: serial must be alphanumeric Latin", ErrValueFormat)
		}
	}
	return nil
}

// RequiredCrypto91 ensures crypto id (91) is exactly 4 alnum chars (Latin letters or digits).
func RequiredCrypto91(value string) error {
	if err := MinLength(value, 4); err != nil {
		return err
	}
	if err := MaxLength(value, 4); err != nil {
		return err
	}
	for _, r := range value {
		if !unicode.IsDigit(r) && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return fmt.Errorf("%w: crypto91 must be alphanumeric Latin", ErrValueFormat)
		}
	}
	return nil
}

// RequiredCrypto92 ensures crypto value (92) is exactly 44 printable characters.
// The spec allows digits, letters and special printable symbols — here we enforce printable ASCII range 32..126.
func RequiredCrypto92(value string) error {
	if err := MinLength(value, 44); err != nil {
		return err
	}
	if err := MaxLength(value, 44); err != nil {
		return err
	}
	for _, r := range value {
		if r < 32 || r > 126 {
			return fmt.Errorf("%w: crypto92 contains non-printable characters", ErrValueFormat)
		}
	}
	return nil
}
