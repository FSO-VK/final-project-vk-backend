package httputil

import (
	"context"
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	"github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
)

const (
	// MsgUnauthorized is a err message for unauthorized user.
	MsgUnauthorized api.ErrorType = "User is not authorized"
	// SessionCookieKey is a key for session cookie.
	SessionCookieKey string = "session_id"
)

// AuthPrincipal contains user id and authorization status for handlers to use.
type AuthPrincipal struct {
	UserID       string
	IsAuthorized bool
}

type ctxKey string

const authUserValueKey ctxKey = "httputil_auth_principal"

// SetAuthToCtx returns a new *http.Request with AuthPrincipal added to the context.
func SetAuthToCtx(r *http.Request, p *AuthPrincipal) *http.Request {
	ctx := context.WithValue(r.Context(), authUserValueKey, p)
	return r.WithContext(ctx)
}

// GetAuthFromCtx extracts AuthPrincipal from the request context.
func GetAuthFromCtx(r *http.Request) (*AuthPrincipal, bool) {
	v := r.Context().Value(authUserValueKey)
	if v == nil {
		return nil, false
	}
	p, ok := v.(*AuthPrincipal)
	return p, ok
}

// AuthMiddleware implements authorization check as a struct.
type AuthMiddleware struct {
	checker    client.AuthChecker
	headerName string
}

// NewAuthMiddleware constructor for AuthMiddleware.
func NewAuthMiddleware(checker client.AuthChecker, headerName string) *AuthMiddleware {
	return &AuthMiddleware{checker: checker, headerName: headerName}
}

// AuthMiddlewareWrapper wraps http.Handler with authorization check.
func (m *AuthMiddleware) AuthMiddlewareWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sid string
		if c, err := r.Cookie(SessionCookieKey); err == nil {
			sid = c.Value
		}

		if sid == "" {
			_ = NetHTTPWriteJSON(w, &api.Response[struct{}]{
				StatusCode: http.StatusUnauthorized,
				Body:       struct{}{},
				Error:      MsgUnauthorized,
			})
			return
		}

		resp, err := m.checker.CheckAuth(&client.Request{SessionID: sid})
		if err != nil {
			_ = NetHTTPWriteJSON(w, &api.Response[struct{}]{
				StatusCode: http.StatusServiceUnavailable,
				Body:       struct{}{},
				Error:      api.MsgServerError,
			})
			return
		}

		if !resp.IsAuthorized {
			_ = NetHTTPWriteJSON(w, &api.Response[struct{}]{
				StatusCode: http.StatusForbidden,
				Body:       struct{}{},
				Error:      MsgUnauthorized,
			})
			return
		}

		r = SetAuthToCtx(r, &AuthPrincipal{
			UserID:       resp.UserID,
			IsAuthorized: resp.IsAuthorized,
		})

		next.ServeHTTP(w, r)
	})
}
