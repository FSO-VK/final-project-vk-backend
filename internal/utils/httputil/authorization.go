package httputil

import (
	"context"
	"errors"
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
	"github.com/gin-gonic/gin"
)

const (
	// SessionCookieKey is a key for session cookie.
	SessionCookieKey string = "session_id"
)

var (
	ErrAuthNotFound    = errors.New("authentication data not found in context")
	ErrInvalidAuthType = errors.New("invalid authentication data type in context")
)

func GetAuthContextKey() interface{} {
	return authUserValueKey
}

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
		status, body, authStatus := m.checkAuth(r)
		if authStatus == nil {
			w.WriteHeader(status)
			_ = NetHTTPWriteJSON(w, body)
			return
		}

		r = SetAuthToCtx(r, authStatus)
		next.ServeHTTP(w, r)
	})
}

// Middleware wraps gin with authorization check.
func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		status, body, authStatus := m.checkAuth(c.Request)
		if authStatus == nil {
			c.AbortWithStatusJSON(status, body)
			return
		}

		r := SetAuthToCtx(c.Request, authStatus)
		c.Request = r
		c.Next()
	}
}

func (m *AuthMiddleware) checkAuth(r *http.Request) (int, *api.Response[any], *AuthStatus) {
	var sid string
	if c, err := r.Cookie(SessionCookieKey); err == nil {
		sid = c.Value
	}

	if sid == "" {
		resp := &api.Response[any]{
			StatusCode: http.StatusUnauthorized,
			Body:       struct{}{},
			Error:      api.MsgUnauthorized,
		}
		return http.StatusUnauthorized, resp, nil
	}

	authResp, err := m.checker.CheckAuth(&auth.Request{SessionID: sid})
	if err != nil {
		resp := &api.Response[any]{
			StatusCode: http.StatusServiceUnavailable,
			Body:       ResponseUnauthorized{SessionID: ""},
			Error:      api.MsgServerError,
		}
		return http.StatusServiceUnavailable, resp, nil
	}

	if !authResp.IsAuthorized {
		resp := &api.Response[any]{
			StatusCode: http.StatusForbidden,
			Body:       struct{}{},
			Error:      api.MsgUnauthorized,
		}
		return http.StatusForbidden, resp, nil
	}

	return http.StatusOK, nil, &AuthStatus{
		UserID:       authResp.UserID,
		IsAuthorized: authResp.IsAuthorized,
	}
}
