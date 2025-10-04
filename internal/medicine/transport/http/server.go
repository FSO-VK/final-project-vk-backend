package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type HTTPServer struct {
	addr   string
	srv    *http.Server
	logger *logrus.Entry
}

func NewHTTPServer(addr string, l *logrus.Entry) *HTTPServer {
	return &HTTPServer{
		addr: addr,
		srv: &http.Server{
			Addr:              addr,
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
	s.logger.Infof("Server started on %s", s.addr)
	return s.srv.ListenAndServe()
}
