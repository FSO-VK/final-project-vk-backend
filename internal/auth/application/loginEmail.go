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
	Email    string
	Password string
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

	sess := session.NewSession(cred.ID)
	err = s.sessionRepo.Create(ctx, sess)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &LoginByEmailResult{
		UserID:    cred.ID.String(),
		SessionID: sess.ID.String(),
		ExpiresAt: sess.ExpiresAt,
	}, nil
}
