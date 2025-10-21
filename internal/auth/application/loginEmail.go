package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/auth/domain/credential"
	"github.com/FSO-VK/final-project-vk-backend/internal/auth/domain/session"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/password"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrNotEmailCredentials = errors.New("credentials type is not email")
)

type LoginByEmail interface {
	Execute(ctx context.Context, login *LoginByEmailCommand) (*LoginByEmailResult, error)
}

// LoginByEmailCommand represents the command to login by email.
type LoginByEmailCommand struct {
	// Optional. If not provided, new session will be created.
	CurrentDeviceSessionID string
	Email                  string `validate:"required,email"`
	Password               string `validate:"required"`
}

// LoginByEmailResult represents the result of a login by email operation.
type LoginByEmailResult struct {
	UserID    string
	SessionID string
	ExpiresAt time.Time
}

type LoginByEmailService struct {
	credentialRepo credential.CredentialRepository
	sessionRepo    session.SessionRepository
	valid          validator.Validator
	passwordHasher password.PasswordHasher
}

func NewLoginByEmailService(
	credentialRepo credential.CredentialRepository,
	sessionRepo session.SessionRepository,
	valid validator.Validator,
	passwordHasher password.PasswordHasher,
) *LoginByEmailService {
	return &LoginByEmailService{
		credentialRepo: credentialRepo,
		sessionRepo:    sessionRepo,
		valid:          valid,
		passwordHasher: passwordHasher,
	}
}

func (s *LoginByEmailService) Execute(
	ctx context.Context,
	loginCmd *LoginByEmailCommand,
) (*LoginByEmailResult, error) {
	err := s.valid.ValidateStruct(loginCmd)
	if err != nil {
		err = errors.Join(err, ErrInvalidCredentials)
		return nil, fmt.Errorf("invalid login command: %w", err)
	}

	cred, err := s.credentialRepo.FindByEmail(ctx, loginCmd.Email)
	if errors.Is(err, credential.ErrNoCredentialFound) {
		return nil, ErrInvalidCredentials
	} else if err != nil {
		return nil, fmt.Errorf("failed to find credential by email: %w", err)
	}

	if !cred.IsTypeEmail() {
		return nil, ErrNotEmailCredentials
	}

	isCorrectPassword := s.passwordHasher.Compare(
		loginCmd.Password,
		cred.Secret.GetSecret(),
	)
	if !isCorrectPassword {
		return nil, ErrInvalidCredentials
	}

	// Check if there is a session for the current device with the same credential.
	currentSessionID, err := uuid.Parse(loginCmd.CurrentDeviceSessionID)
	if err == nil {
		currentSession, err := s.sessionRepo.GetByID(ctx, currentSessionID)
		if err == nil {
			return &LoginByEmailResult{
				UserID:    cred.ID.String(),
				SessionID: currentSession.ID.String(),
				ExpiresAt: currentSession.ExpiresAt,
			}, nil
		}
	}

	newSession := session.NewSession(cred.ID)
	err = s.sessionRepo.Create(ctx, newSession)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &LoginByEmailResult{
		UserID:    cred.ID.String(),
		SessionID: newSession.ID.String(),
		ExpiresAt: newSession.ExpiresAt,
	}, nil
}
