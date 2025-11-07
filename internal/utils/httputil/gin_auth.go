package httputil

import (
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
	"github.com/gin-gonic/gin"
)

type GinAuthMiddleware struct {
	checker auth.AuthChecker
}

func NewGinAuthMiddleware(checker auth.AuthChecker) *GinAuthMiddleware {
	return &GinAuthMiddleware{checker: checker}
}

func (m *GinAuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var sid string
		if cookie, err := c.Request.Cookie(SessionCookieKey); err == nil {
			sid = cookie.Value
		}

		if sid == "" {
			_ = NetHTTPWriteJSON(c.Writer, &api.Response[struct{}]{
				StatusCode: http.StatusUnauthorized,
				Body:       struct{}{},
				Error:      api.MsgUnauthorized,
			})
			c.Abort()
			return
		}

		resp, err := m.checker.CheckAuth(&auth.Request{SessionID: sid})
		if err != nil {
			_ = NetHTTPWriteJSON(c.Writer, &api.Response[ResponseUnauthorized]{
				StatusCode: http.StatusServiceUnavailable,
				Body:       ResponseUnauthorized{SessionID: ""},
				Error:      api.MsgServerError,
			})
			c.Abort()
			return
		}

		if !resp.IsAuthorized {
			_ = NetHTTPWriteJSON(c.Writer, &api.Response[struct{}]{
				StatusCode: http.StatusForbidden,
				Body:       struct{}{},
				Error:      api.MsgUnauthorized,
			})
			c.Abort()
			return
		}

		r := SetAuthToCtx(c.Request, &AuthStatus{
			UserID:       resp.UserID,
			IsAuthorized: resp.IsAuthorized,
		})
		c.Request = r

		c.Next()
	}
}
