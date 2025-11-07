package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host string
	Port string
}

// Unexported global variable.
//
//nolint:gochecknoglobals
var defaultServerConfig = &ServerConfig{
	Host: "",
	Port: "8000",
}

// Address returns the full server address in host:port format.
func (s *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

// MedicationHTTPServer represents an HTTP server with configuration, logger and underlying http.Server.
type MedicationHTTPServer struct {
	config *ServerConfig
	srv    *http.Server
	logger *logrus.Entry
}

// NewGINServer creates a new HTTP server instance with the provided configuration and logger.
func NewGINServer(conf *ServerConfig, l *logrus.Entry) *MedicationHTTPServer {
	if conf == nil {
		conf = defaultServerConfig
	}
	return &MedicationHTTPServer{
		config: conf,
		srv: &http.Server{
			Addr:                         conf.Address(),
			Handler:                      nil,
			ReadHeaderTimeout:            5 * time.Second,
			DisableGeneralOptionsHandler: false,
			TLSConfig:                    nil,
			ReadTimeout:                  0,
			WriteTimeout:                 0,
			IdleTimeout:                  0,
			MaxHeaderBytes:               0,
			TLSNextProto:                 nil,
			ConnState:                    nil,
			ErrorLog:                     nil,
			BaseContext:                  nil,
			ConnContext:                  nil,
			HTTP2:                        nil,
			Protocols:                    nil,
		},
		logger: l,
	}
}

// Router sets the HTTP router for the server.
func (s *MedicationHTTPServer) Router(router *gin.Engine) {
	s.srv.Handler = router
}

// Shutdown gracefully shuts down the server.
func (s *MedicationHTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// ListenAndServe starts the HTTP server.
func (s *MedicationHTTPServer) ListenAndServe() error {
	s.logger.Infof("Server started on %s", s.config.Address())
	return s.srv.ListenAndServe()
}
