package httputil

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

type PanicRecoveryMiddleware struct {
	log *log.Logger
}

func NewPanicRecoveryMiddleware() *PanicRecoveryMiddleware {
	// TODO: refactor
	return &PanicRecoveryMiddleware{log: log.Default()}
}

func (p *PanicRecoveryMiddleware) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recovery(p)
		handler.ServeHTTP(w, r)
	})
}

func (p *PanicRecoveryMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		recovery(p)
		c.Next()
	}
}

func recovery(p *PanicRecoveryMiddleware) {
	if err := recover(); err != nil {
		p.log.Printf("Panic recovered: %v", err)

		buf := make([]byte, 1024)
		n := runtime.Stack(buf, false)
		for n == len(buf) {
			buf = make([]byte, len(buf)*2)
			n = runtime.Stack(buf, false)
		}

		fmt.Printf("Stack trace:\n%s\n", buf[:n])

	}
}
