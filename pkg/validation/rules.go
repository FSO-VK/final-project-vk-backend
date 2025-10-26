// Package validation contains validation rules.
package validation

import (
	"errors"
	"fmt"
)

var (
	ErrValueRequired = errors.New("value can't be empty")
	ErrValueShort    = errors.New("value is too short")
	ErrValueLong     = errors.New("value is too long")
	ErrValueNegative = errors.New("value can't be negative")
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
