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
	ErrInvalidEmail     = errors.New("invalid email")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUserAlreadyExist = errors.New("user already exist")
)

type Registration interface {
	Execute(ctx context.Context, login *RegistrationCommand) (*RegistrationResult, error)
}

// RegistrationCommand represents the command to register a new user.
type RegistrationCommand struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=64"`
}

// RegistrationResult represents the result of a registration operation.
type RegistrationResult struct {
	UserID    string
	SessionID string
	ExpiresAt time.Time
}

type RegistrationService struct {
	credentialRepo credential.CredentialRepository
	sessionRepo    session.SessionRepository
	valid          validator.Validator
	passwordHasher password.PasswordHasher
}

func NewRegistrationService(
	credentialRepo credential.CredentialRepository,
	sessionRepo session.SessionRepository,
	valid validator.Validator,
	passwordHasher password.PasswordHasher,
) *RegistrationService {
	return &RegistrationService{
		credentialRepo: credentialRepo,
		sessionRepo:    sessionRepo,
		valid:          valid,
		passwordHasher: passwordHasher,
	}
}

func (s *RegistrationService) Execute(
	ctx context.Context,
	registrationCmd *RegistrationCommand,
) (*RegistrationResult, error) {
	err := s.valid.ValidateStruct(registrationCmd)
	if err != nil {
		return nil, fmt.Errorf("invalid login command: %w", err)
	}

	Type := credential.TypeEmail
	_, err1 := s.credentialRepo.FindByEmail(ctx, registrationCmd.Email)
	ID := uuid.New()
	if err1 == nil {
		return nil, ErrUserAlreadyExist
	}
	for {
		_, err1 := s.credentialRepo.FindByID(ctx, ID)
		if err1 != nil {
			break
		}
		ID = uuid.New()
	}
	Password, err := credential.NewSecretPassword(registrationCmd.Password, s.passwordHasher)
	if err != nil {
		return nil, ErrInvalidPassword
	}
	user := credential.NewCredential(
		ID,
		Type,
		registrationCmd.Email,
		Password,
		time.Now(),
	)
	if !user.IsTypeEmail() {
		return nil, ErrNotEmailCredentials
	}
	err = s.credentialRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to find credential by email: %w", err)
	}

	sess := session.NewSession(user.ID)
	err = s.sessionRepo.Create(ctx, sess)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &RegistrationResult{
		UserID:    user.ID.String(),
		SessionID: sess.ID.String(),
		ExpiresAt: sess.ExpiresAt,
	}, nil
}
