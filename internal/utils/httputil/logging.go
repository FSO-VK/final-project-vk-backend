package httputil

import (
	"bufio"
	"net"
	"net/http"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/utils/logcon"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// statusResponse is a wrapper for http.ResponseWriter to capture status code.
type statusResponse struct {
	http.ResponseWriter

	statusCode int
}

// WriteHeader overrides http.ResponseWriter.WriteHeader to capture status code.
func (s *statusResponse) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

// Hijack checks if the underlying response implements http.Hijacker interface.
func (s *statusResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := s.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hj.Hijack()
}

// LoggingMiddleware is a middleware that logs incoming requests.
type LoggingMiddleware struct {
	log *logrus.Entry
}

// NewLoggingMiddleware creates a new LoggingMiddleware.
func NewLoggingMiddleware(log *logrus.Entry) *LoggingMiddleware {
	return &LoggingMiddleware{log: log}
}

// MiddlewareNetHTTP is a middleware method for net/http compatible frameworks.
func (m *LoggingMiddleware) MiddlewareNetHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := m.log.WithFields(logrus.Fields{
			"method":     r.Method,
			"url":        r.URL.String(),
			"request_id": uuid.New().String(),
		})
		log.Info("Incoming request")

		ctx := logcon.WithContext(r.Context(), log)
		rw := &statusResponse{ResponseWriter: w}

		start := time.Now()
		next.ServeHTTP(rw, r.WithContext(ctx))
		end := time.Since(start)

		log.WithFields(logrus.Fields{
			"status": rw.statusCode,
			"timing": end.String(),
		}).Info("Request handled")
	})
}
