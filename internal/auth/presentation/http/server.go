package http

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type ServerConfig struct {
	Host string
	Port string
}

func (s *ServerConfig) Address() string {
	return s.Host + ":" + s.Port
}

type ServerHTTP struct {
	conf *ServerConfig
	srv  *fasthttp.Server
	// logger must at least implement fasthttp.Logger's
	// method Printf()
	logger *logrus.Entry
}

func NewServerHTTP(conf ServerConfig, router *Router, logger *logrus.Entry) *ServerHTTP {
	return &ServerHTTP{
		conf: &conf,
		srv: &fasthttp.Server{
			Logger:  logger,
			Handler: RequestInfoMiddleware(router.GetRouter(), logger),
		},
		logger: logger,
	}
}

func (s *ServerHTTP) ListenAndServe() error {
	s.logger.Printf("Server started on %s", s.conf.Address())
	return s.srv.ListenAndServe(s.conf.Address())
}

// Shutdown gracefully shutdowns the server.
func (s *ServerHTTP) Shutdown() error {
	return s.srv.Shutdown()
}
