// Package prober provides a simple way to define application probes.
package prober

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/pkg/prober/edge"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host string
	Port string
}

// Address returns the full server address in host:port format.
func (s *Config) Address() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

// Prober is a wrapper for http server.
type Prober struct {
	config *Config
	srv    *http.Server
	logger *logrus.Entry
}

// New makes new Prober.
func New(appEdge *edge.Edge, conf *Config, l *logrus.Entry) Prober {
	handlers := newHandlers(appEdge, l)

	router := http.NewServeMux()
	router.HandleFunc("/healthz", handlers.Health)
	router.HandleFunc("/readyz", handlers.Ready)

	srv := &http.Server{
		Addr:                         conf.Address(),
		Handler:                      router,
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
	}

	return Prober{
		config: conf,
		srv:    srv,
		logger: l,
	}
}

// ListenAndServe raises prober server.
func (p *Prober) ListenAndServe() error {
	p.logger.Infof("Prober started on %s", p.config.Address())
	return p.srv.ListenAndServe()
}

// Shutdown gracefully shutdowns the prober.
func (s *Prober) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
