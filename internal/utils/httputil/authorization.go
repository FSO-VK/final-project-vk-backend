package httputil

import (
	"context"
	"errors"
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
)

const (
	// SessionCookieKey is a key for session cookie.
	SessionCookieKey string = "session_id"
)

var (
	ErrAuthNotFound    = errors.New("authentication data not found in context")
	ErrInvalidAuthType = errors.New("invalid authentication data type in context")
)

// AuthStatus contains user id and authorization status for handlers to use.
type AuthStatus struct {
	UserID       string
	IsAuthorized bool
}

// ResponseUnauthorized is a response for unauthorized requests.
type ResponseUnauthorized struct {
	SessionID string `json:"sessionId"`
}

type ctxKey string

const authUserValueKey ctxKey = "httputil_auth_principal"

// SetAuthToCtx returns a new *http.Request with AuthStatus added to the context.
func SetAuthToCtx(r *http.Request, p *AuthStatus) *http.Request {
	ctx := context.WithValue(r.Context(), authUserValueKey, p)
	return r.WithContext(ctx)
}

func GetAuthFromCtx(r *http.Request) (*AuthStatus, error) {
	v := r.Context().Value(authUserValueKey)
	if v == nil {
		return nil, ErrAuthNotFound
	}
	p, ok := v.(*AuthStatus)
	if !ok {
		return nil, ErrInvalidAuthType
	}
	return p, nil
}

// AuthMiddleware implements authorization check as a struct.
type AuthMiddleware struct {
	checker auth.AuthChecker
}

// NewAuthMiddleware constructor for AuthMiddleware.
func NewAuthMiddleware(checker auth.AuthChecker) *AuthMiddleware {
	return &AuthMiddleware{checker: checker}
}

// AuthMiddlewareWrapper wraps http.Handler with authorization check.
func (m *AuthMiddleware) AuthMiddlewareWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sid string
		if c, err := r.Cookie(SessionCookieKey); err == nil {
			sid = c.Value
		}

		if sid == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_ = NetHTTPWriteJSON(w, &api.Response[struct{}]{
				StatusCode: http.StatusUnauthorized,
				Body:       struct{}{},
				Error:      api.MsgUnauthorized,
			})
			return
		}

		resp, err := m.checker.CheckAuth(&auth.Request{SessionID: sid})
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = NetHTTPWriteJSON(w, &api.Response[ResponseUnauthorized]{
				StatusCode: http.StatusServiceUnavailable,
				Body:       ResponseUnauthorized{SessionID: ""},
				Error:      api.MsgServerError,
			})
			return
		}

		if !resp.IsAuthorized {
			w.WriteHeader(http.StatusForbidden)
			_ = NetHTTPWriteJSON(w, &api.Response[struct{}]{
				StatusCode: http.StatusForbidden,
				Body:       struct{}{},
				Error:      api.MsgUnauthorized,
			})
			return
		}

		r = SetAuthToCtx(r, &AuthStatus{
			UserID:       resp.UserID,
			IsAuthorized: resp.IsAuthorized,
		})

		next.ServeHTTP(w, r)
	})
}
