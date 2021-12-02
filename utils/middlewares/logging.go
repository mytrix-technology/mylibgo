package middlewares

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

func NewEndpointLoggingMiddleware(logger log.Logger, endpointMethod string) endpoint.Middleware {
	return MakeEndpointLoggingMiddleware(logger, endpointMethod)
}
