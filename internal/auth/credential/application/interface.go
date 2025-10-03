package application

import (
	"context"
	"fmt"
	"time"
)

var (
	ErrNoMatchPasswords = fmt.Errorf("passwords do not match")
)

type CredentialServiceProvider interface {
	LoginByEmail(ctx context.Context, login *LoginEmailRequest) (*LoginEmailResponse, error)
	Logout(ctx context.Context, logout *LogoutRequest) (*LogoutResponse, error)
}

type LoginEmailRequest struct {
	Email           string `validate:"required,email"`
	Password        string `validate:"required"`
	PasswordConfirm string `validate:"required"`
}

type LoginEmailResponse struct {
	SessionID string
	ExpiresAt time.Time
}

type LogoutRequest struct {
	SessionID string `validate:"required"`
}

type LogoutResponse struct {
	SessionID string
}
