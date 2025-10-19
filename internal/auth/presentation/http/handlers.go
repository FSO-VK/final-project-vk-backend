package http

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/auth/application"
	httph "github.com/FSO-VK/final-project-vk-backend/internal/transport/http"
	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	SessionCookieKey = "session_id"
)

type AuthHandlers struct {
	app    *application.AuthApplication
	logger *logrus.Entry
}

func NewAuthHandlers(
	app *application.AuthApplication,
	logger *logrus.Entry,
) *AuthHandlers {
	return &AuthHandlers{
		app:    app,
		logger: logger,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID string `json:"userId"`
}

func (h *AuthHandlers) Login(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	if len(body) == 0 {
		h.logger.Error("Empty request body")

		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgNoBody,
		})

		return
	}

	var req LoginRequest
	err := json.Unmarshal(body, &req)
	if err != nil {
		h.logger.WithError(err).Errorf("Failed to read request body: %v", err)

		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgNoBody,
		})

		return
	}

	sessionID := ctx.Request.Header.Cookie(SessionCookieKey)

	serviceRequest := &application.LoginByEmailCommand{
		CurrentDeviceSessionID: string(sessionID),
		Email:                  req.Email,
		Password:               req.Password,
	}

	serviceResult, err := h.app.LoginByEmail.Execute(ctx, serviceRequest)
	if errors.Is(err, application.ErrInvalidCredentials) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusUnauthorized,
			Body:       struct{}{},
			Error:      MsgWrongCredentials,
		})

		return
	} else if err != nil {
		h.logger.WithError(err).Error("Failed to login by email")

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusInternalServerError,
			Body:       struct{}{},
			Error:      api.MsgServerError,
		})

		return
	}

	response := &LoginResponse{
		UserID: serviceResult.UserID,
	}

	err = setSessionCookie(
		ctx,
		SessionCookieKey,
		serviceResult.SessionID,
		serviceResult.ExpiresAt,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to set session cookie")

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgSetCookieFail,
		})

		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	_ = httph.FastHTTPWriteJSON(ctx, &api.Response[*LoginResponse]{
		StatusCode: fasthttp.StatusOK,
		Body:       response,
		Error:      "",
	})
}

func (h *AuthHandlers) Logout(ctx *fasthttp.RequestCtx) {
	sessionID := ctx.Request.Header.Cookie(SessionCookieKey)

	if len(sessionID) == 0 {
		ctx.SetStatusCode(fasthttp.StatusOK)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusOK,
			Body:       struct{}{},
			Error:      "",
		})

		return
	}

	serviceRequest := &application.LogoutCommand{
		SessionID: string(sessionID),
	}

	_, err := h.app.Logout.Execute(ctx, serviceRequest)
	if err != nil {
		h.logger.WithError(err).Error("Failed to logout")

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgLogoutFailed,
		})

		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
		StatusCode: fasthttp.StatusOK,
		Body:       struct{}{},
		Error:      "",
	})
}

type CheckAuthResponse struct {
	UserID string `json:"userId"`
}

type CheckAuthResponseFail struct {
	SessionID string `json:"sessionId"`
}

func (h *AuthHandlers) CheckAuth(ctx *fasthttp.RequestCtx) {
	sessionID := ctx.Request.Header.Cookie(SessionCookieKey)

	if len(sessionID) == 0 {
		h.logger.Errorf("No %s in cookie", SessionCookieKey)

		ctx.SetStatusCode(fasthttp.StatusUnauthorized)

		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[*CheckAuthResponseFail]{
			StatusCode: fasthttp.StatusUnauthorized,
			Body: &CheckAuthResponseFail{
				SessionID: "",
			},
			Error: MsgUnauthorized,
		})

		return
	}

	serviceRequest := &application.CheckAuthCommand{
		SessionID: string(sessionID),
	}

	serviceResult, err := h.app.CheckAuth.Execute(ctx, serviceRequest)
	if err != nil {
		h.logger.WithError(err).Error("Failed to check auth")

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[*CheckAuthResponseFail]{
			StatusCode: fasthttp.StatusInternalServerError,
			Body: &CheckAuthResponseFail{
				SessionID: serviceRequest.SessionID,
			},
			Error: api.MsgServerError,
		})

		return
	}

	if !serviceResult.IsAuthenticated {
		h.logger.WithError(err).Errorf("User is not authenticated")

		ctx.SetStatusCode(fasthttp.StatusUnauthorized)

		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[*CheckAuthResponseFail]{
			StatusCode: fasthttp.StatusUnauthorized,
			Body: &CheckAuthResponseFail{
				SessionID: serviceRequest.SessionID,
			},
			Error: MsgUnauthorized,
		})

		return
	}

	err = setSessionCookie(
		ctx,
		SessionCookieKey,
		serviceResult.SessionID,
		serviceResult.ExpiresAt,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to set session cookie")

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgSetCookieFail,
		})

		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	_ = httph.FastHTTPWriteJSON(ctx, &api.Response[*CheckAuthResponse]{
		StatusCode: fasthttp.StatusOK,
		Body: &CheckAuthResponse{
			UserID: serviceResult.SessionID,
		},
		Error: "",
	})
}

var ErrNoCookie = errors.New("cookie is nil")

func setSessionCookie(
	ctx *fasthttp.RequestCtx,
	seesionID string,
	value string,
	expiration time.Time,
) error {
	c := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(c)

	if c == nil {
		return ErrNoCookie
	}

	c.SetKey(seesionID)
	c.SetValue(value)
	c.SetExpire(expiration)
	c.SetHTTPOnly(true)
	c.SetSecure(true)

	ctx.Response.Header.SetCookie(c)
	return nil
}

type RegistrationByEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationByEmailResponse struct {
	UserID string `json:"userId"`
}

func (h *AuthHandlers) RegistrationByEmail(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	if len(body) == 0 {
		h.logger.Error("Empty request body")

		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgNoBody,
		})

		return
	}

	var req RegistrationByEmailRequest
	err := json.Unmarshal(body, &req)
	if err != nil {
		h.logger.WithError(err).Errorf("Failed to read request body: %v", err)

		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgNoBody,
		})

		return
	}

	serviceRequest := &application.RegistrationCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	serviceResult, err := h.app.Registration.Execute(ctx, serviceRequest)
	if err != nil {
		h.logger.WithError(err).Error("Failed to register user")

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		errorOut := MsgInvalidEmail
		statusCode := fasthttp.StatusInternalServerError
		if errors.Is(err, application.ErrInvalidPassword) {
			errorOut = MsgInvalidPassword
			statusCode = fasthttp.StatusBadRequest
		}
		if errors.Is(err, application.ErrUserAlreadyExist) {
			errorOut = MsgUserAlreadyExist
			statusCode = fasthttp.StatusBadRequest
		}
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: statusCode,
			Body:       struct{}{},
			Error:      errorOut,
		})

		return
	}

	response := &RegistrationByEmailResponse{
		UserID: serviceResult.UserID,
	}

	err = setSessionCookie(
		ctx,
		SessionCookieKey,
		serviceResult.SessionID,
		serviceResult.ExpiresAt,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to set session cookie")

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgSetCookieFail,
		})

		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	_ = httph.FastHTTPWriteJSON(ctx, &api.Response[*RegistrationByEmailResponse]{
		StatusCode: fasthttp.StatusOK,
		Body:       response,
		Error:      "",
	})
}
