package httputil

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"
)

type GINPanicRecoveryMiddleware struct {
	log *log.Logger
}

func NewGINPanicRecoveryMiddleware() *GINPanicRecoveryMiddleware {
	return &GINPanicRecoveryMiddleware{log: log.Default()}
}

func (p *GINPanicRecoveryMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				p.log.Printf("Panic recovered: %v", err)

				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)
				for n == len(buf) {
					buf = make([]byte, len(buf)*2)
					n = runtime.Stack(buf, false)
				}

				fmt.Printf("Stack trace:\n%s\n", buf[:n])

				c.AbortWithStatusJSON(500, gin.H{
					"error": "internal server error",
				})
			}
		}()

		c.Next()
	}
}
