package validation_test

import (
	"testing"

	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

func TestEAN13(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		value   string
		wantErr bool
	}{
		{
			name:    "Should return nil error for valid EAN13",
			value:   "4605077016860",
			wantErr: false,
		},
		{
			name:    "Should return error if syntax is ok but check digit is not",
			value:   "460507701686",
			wantErr: true,
		},
		{
			name:    "Should return error if EAN13 contains non-digit",
			value:   "46a4jk7016860",
			wantErr: true,
		},
		{
			name:    "Should return error if EAN13 is not 13 symbols long",
			value:   "12345",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := validation.EAN13(tt.value)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("EAN13() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("EAN13() succeeded unexpectedly")
			}
		})
	}
}

func TestFixedLength(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		value   string
		length  int
		wantErr bool
	}{
		{
			name:    "Should return nil error if string has Unicode",
			value:   "русский text",
			length:  12,
			wantErr: false,
		},
		{
			name:    "Should return error if string longer then desired",
			value:   "русский text",
			length:  5,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := validation.FixedLength(tt.value, tt.length)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("FixedLength() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("FixedLength() succeeded unexpectedly")
			}
		})
	}
}
