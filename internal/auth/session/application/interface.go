package application

import (
	"context"
	"time"
)

type SessionService interface {
	CheckAuth(ctx context.Context, checkReq *CheckAuthRequest) (*CheckAuthResponse, error)
}

type CheckAuthRequest struct {
	SessionID string `validate:"required,uuid"`
}

type CheckAuthResponse struct {
	SessionID    string
	IsAuthorized bool
	ExpiresAt    time.Time
}
