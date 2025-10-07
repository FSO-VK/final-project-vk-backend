package http

import (
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
	switch string(ctx.Path()) {
	case "api/v1/auth/login":
		r.withMethod(r.handlers.Login, MethodPost)(ctx)
	case "api/v1/auth/logout":
		r.withMethod(r.handlers.Logout, MethodPost)(ctx)
	case "api/v1/auth/check":
		r.withMethod(r.handlers.CheckAuth, MethodGet)(ctx)
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
