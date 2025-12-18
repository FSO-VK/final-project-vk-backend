package http

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func RequestInfoMiddleware(
	next fasthttp.RequestHandler,
	logger *logrus.Entry,
) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		l := logger.WithFields(logrus.Fields{
			"method": string(ctx.Method()),
			"uri":    string(ctx.RequestURI()),
		})
		l.Info("Handle request")

		start := time.Now()

		next(ctx)

		end := time.Since(start)

		l.WithFields(logrus.Fields{
			"status": ctx.Response.StatusCode(),
			"timing": end.String(),
		}).Info("Request handled")
	}
}
