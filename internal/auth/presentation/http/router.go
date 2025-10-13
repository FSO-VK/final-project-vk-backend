package http

import (
	"strings"

	httph "github.com/FSO-VK/final-project-vk-backend/internal/transport/http"
	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	"github.com/valyala/fasthttp"
)

type Method string

const (
	MethodGet    Method = "GET"
	MethodPost   Method = "POST"
	MethodPut    Method = "PUT"
	MethodDelete Method = "DELETE"
	MethodPatch  Method = "PATCH"
)

type Router struct {
	handlers *AuthHandlers
}

func NewRouter(handlers *AuthHandlers) *Router {
	return &Router{
		handlers: handlers,
	}
}

func (r *Router) GetRouter() fasthttp.RequestHandler {
	return r.router
}

func (r *Router) router(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	method := strings.ToUpper(string(ctx.Method()))

	switch path {
	case "/session":
		switch method {
		case string(MethodPost):
			r.withMethod(r.handlers.Login, MethodPost)(ctx)
		case string(MethodGet):
			r.withMethod(r.handlers.CheckAuth, MethodGet)(ctx)
		case string(MethodDelete):
			r.withMethod(r.handlers.Logout, MethodDelete)(ctx)
		default:
			r.handlerMethodNotAllowed(ctx)
		}
	case "/user":
		switch method {
		case string(MethodPost):
			r.withMethod(r.handlers.RegistrationByEmail, MethodPost)(ctx)
		default:
			r.handlerMethodNotAllowed(ctx)
		}
	default:
		r.handlerNotFound(ctx)
	}
}

func (r *Router) withMethod(
	handler func(ctx *fasthttp.RequestCtx),
	methods ...Method,
) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		for _, method := range methods {
			if Method(ctx.Method()) == method {
				handler(ctx)
				return
			}
		}
		r.handlerMethodNotAllowed(ctx)
	}
}

func (r *Router) handlerMethodNotAllowed(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
		StatusCode: fasthttp.StatusMethodNotAllowed,
		Body:       struct{}{},
		Error:      api.MsgMethodNotAllowed,
	})
}

func (r *Router) handlerNotFound(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	_ = httph.FastHTTPWriteJSON(ctx, &api.Response[struct{}]{
		StatusCode: fasthttp.StatusNotFound,
		Body:       struct{}{},
		Error:      api.MsgNotFound,
	})
}
