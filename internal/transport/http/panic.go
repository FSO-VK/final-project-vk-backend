package http

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
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
		defer func() {
			if err := recover(); err != nil {
				p.log.Printf("Panic recovered: %v", err)

				buf := make([]byte, 1024)

				n := runtime.Stack(buf, false)
				for n == len(buf) {
					buf = make([]byte, len(buf)*2)
					n = runtime.Stack(buf, false)
				}
				fmt.Printf("Stack trace: %s\n", buf[:n])
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
