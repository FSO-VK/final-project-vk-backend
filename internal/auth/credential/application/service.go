package application

import (
	"context"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/utils/password"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/google/uuid"

	credomain "github.com/FSO-VK/final-project-vk-backend/internal/auth/credential/domain"
	sesdomain "github.com/FSO-VK/final-project-vk-backend/internal/auth/session/domain"
)

type CredentialService struct {
	crepo          credomain.CredentialRepository
	sessionRepo    sesdomain.SessionRepository
	valid          validator.Validator
	passwordHasher password.PasswordHasher
}

func NewCredentialService(
	crepo credomain.CredentialRepository,
	sessionRepo sesdomain.SessionRepository,
	valid validator.Validator,
	passwordHasher password.PasswordHasher,
) *CredentialService {
	return &CredentialService{
		crepo:          crepo,
		sessionRepo:    sessionRepo,
		valid:          valid,
		passwordHasher: passwordHasher,
	}
}

func (s *CredentialService) LoginByEmail(
	ctx context.Context,
	loginReq *LoginEmailRequest,
) (*LoginEmailResponse, error) {

	err := s.valid.ValidateStruct(loginReq)
	if err != nil {
		return nil, fmt.Errorf("invalid login request: %w", err)
	}

	if loginReq.Password != loginReq.PasswordConfirm {
		return nil, ErrNoMatchPasswords
	}

	password, err := credomain.NewPassword(loginReq.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	cred, err := s.crepo.FindByEmail(ctx, loginReq.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find credential by email: %w", err)
	}

	isActual := s.passwordHasher.Compare(password.String(), cred.Secret)
	if !isActual {
		return nil, fmt.Errorf("in")
	}

	session := sesdomain.NewSession(cred.ID)

	err = s.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &LoginEmailResponse{
		SessionID: session.ID.String(),
		ExpiresAt: session.ExpiresAt,
	}, nil
}

func (s *CredentialService) Logout(
	ctx context.Context,
	logoutReq *LogoutRequest,
) (*LogoutResponse, error) {

	err := s.valid.ValidateStruct(logoutReq)
	if err != nil {
		return nil, fmt.Errorf("invalid logout request: %w", err)
	}

	sessionIDuuid, err := uuid.Parse(logoutReq.SessionID)
	if err != nil {
		return nil, fmt.Errorf("parse session id to uuid: %w", err)
	}

	err = s.sessionRepo.Delete(ctx, sessionIDuuid)
	if err != nil {
		return nil, fmt.Errorf("failed to delete session by id: %w", err)
	}

	return &LogoutResponse{
		SessionID: sessionIDuuid.String(),
	}, nil
}
