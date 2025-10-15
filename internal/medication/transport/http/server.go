package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Host string
	Port string
}

// Unexported global variable.
//
//nolint:gochecknoglobals
var defaultServerConfig = &ServerConfig{
	Host: "",
	Port: "8080",
}

func (s *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

type HTTPServer struct {
	config *ServerConfig
	srv    *http.Server
	logger *logrus.Entry
}

func NewHTTPServer(conf *ServerConfig, l *logrus.Entry) *HTTPServer {
	if conf == nil {
		conf = defaultServerConfig
	}
	return &HTTPServer{
		config: conf,
		srv: &http.Server{
			Addr:              conf.Address(),
			Handler:           nil,
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger: l,
	}
}

func (s *HTTPServer) Router(router *mux.Router) {
	s.srv.Handler = router
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *HTTPServer) ListenAndServe() error {
	s.logger.Infof("Server started on %s", s.config.Address())
	return s.srv.ListenAndServe()
}
