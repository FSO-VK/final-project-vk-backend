// Package validation contains validation rules.
package validation

import (
	"errors"
	"fmt"
	"unicode"
)

var (
	ErrValueRequired    = errors.New("value can't be empty")
	ErrValueShort       = errors.New("value is too short")
	ErrValueLong        = errors.New("value is too long")
	ErrValueNegative    = errors.New("value can't be negative")
	ErrValueFormat      = errors.New("value has invalid format")
	ErrValueFixedLength = errors.New("value must have fixed length")
	ErrNoEAN13          = errors.New("value is no EAN-13")
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
		return fmt.Errorf("%w: can't be longer than %d", ErrValueLong, length)
	}
	return nil
}

func FixedLength(value string, length int) error {
	if len(value) != length {
		return fmt.Errorf("%w: %d", ErrValueFixedLength, length)
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

// GTIN ensures GTIN contains only digits.
func GTIN(value string) error {
	if err := FixedLength(value, 14); err != nil {
		return err
	}
	for _, r := range value {
		if !unicode.IsDigit(r) {
			return fmt.Errorf("%w: gtin must contain only digits", ErrValueFormat)
		}
	}
	return nil
}

// Serial ensures SerialNumber contains only Latin letters and digits.
func Serial(value string) error {
	if err := FixedLength(value, 13); err != nil {
		return err
	}
	for _, r := range value {
		if !unicode.IsDigit(r) && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return fmt.Errorf("%w: serial must be alphanumeric Latin", ErrValueFormat)
		}
	}
	return nil
}

// Crypto91 ensures crypto id (91) contains only Latin letters and digits.
func Crypto91(value string) error {
	if err := FixedLength(value, 4); err != nil {
		return err
	}
	for _, r := range value {
		if !unicode.IsDigit(r) && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return fmt.Errorf("%w: crypto91 must be alphanumeric Latin", ErrValueFormat)
		}
	}
	return nil
}

// Crypto92 ensures crypto value (92) contains only printable ASCII characters.
func Crypto92(value string) error {
	if err := FixedLength(value, 44); err != nil {
		return err
	}
	for _, r := range value {
		if r < 32 || r > 126 {
			return fmt.Errorf("%w: crypto92 contains non-printable characters", ErrValueFormat)
		}
	}
	return nil
}

func EAN13(value string) error {
	err := FixedLength(value, 13)
	if err != nil {
		return ErrNoEAN13
	}

	var ean13 [13]int
	for i, r := range value {
		if !unicode.IsDigit(r) {
			return ErrNoEAN13
		}
		ean13[i] = int(r - '0')
	}

	if !isCorrectEAN13Checksum(ean13) {
		return ErrNoEAN13
	}
	return nil
}

func isCorrectEAN13Checksum(ean13 [13]int) bool {
	sum := 0
	for i, v := range ean13[:12] {
		if i%2 == 0 {
			sum += 3 * v
		} else {
			sum += v
		}
	}

	checkDigit := 10 - (sum % 10)
	return checkDigit == ean13[12]
}
