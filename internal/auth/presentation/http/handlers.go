package http

import (
	"bytes"
	"encoding/json"

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
	var body *bytes.Buffer
	_, err := body.Read(ctx.PostBody())
	if err != nil {
		h.logger.Errorf("Failed to read request body: %v", err)

		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgNoBody,
		})

		return
	}

	var req LoginRequest
	err = json.Unmarshal(body.Bytes(), &req)
	if err != nil {
		h.logger.Errorf("Failed to read request body: %v", err)

		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusBadRequest,
			Body:       struct{}{},
			Error:      api.MsgNoBody,
		})

		return
	}

	serviceRequest := &application.LoginByEmailCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	serviceResult, err := h.app.LoginByEmail.Execute(ctx, serviceRequest)
	if err != nil {
		h.logger.WithError(err).Error("Failed to login by email")

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
			StatusCode: fasthttp.StatusInternalServerError,
			Body:       struct{}{},
			Error:      MsgWrongCredentials,
		})

		return
	}

	response := &LoginResponse{
		UserID: serviceResult.UserID,
	}

	c := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(c)

	c.SetKey(SessionCookieKey)
	c.SetValue(serviceResult.Token)
	c.SetExpire(serviceResult.ExpiresAt)
	c.SetHTTPOnly(true)

	ctx.Response.Header.SetCookie(c)

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

	serviceResponse, err := h.app.CheckAuth.Execute(ctx, serviceRequest)
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

	if !serviceResponse.IsAuthenticated {
		h.logger.Errorf("User is not authenticated")

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

	c := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(c)

	c.SetKey(SessionCookieKey)
	c.SetValue(serviceResponse.SessionID)
	c.SetExpire(serviceResponse.ExpiresAt)
	c.SetHTTPOnly(true)

	ctx.SetStatusCode(fasthttp.StatusOK)
	_ = httph.FastHTTPWriteJSON(ctx, &api.Response[*CheckAuthResponse]{
		StatusCode: fasthttp.StatusOK,
		Body: &CheckAuthResponse{
			UserID: serviceResponse.SessionID,
		},
		Error: "",
	})
}
