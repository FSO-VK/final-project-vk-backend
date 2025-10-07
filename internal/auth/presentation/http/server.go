package http

import "github.com/valyala/fasthttp"

type ServerHTTP struct {
	addr string
	srv  *fasthttp.Server
	// logger must at least implement fasthttp.Logger's
	// method Printf()
	logger fasthttp.Logger
}

func NewServerHTTP(addr string, router *Router, logger fasthttp.Logger) *ServerHTTP {
	return &ServerHTTP{
		addr: addr,
		srv: &fasthttp.Server{
			Logger:  logger,
			Handler: router.GetRouter(),
		},
		logger: logger,
	}
}

func (s *ServerHTTP) ListenAndServe() error {
	return s.srv.ListenAndServe(s.addr)
}

// Shutdown gracefully shutdowns the server.
func (s *ServerHTTP) Shutdown() error {
	return s.srv.Shutdown()
}
