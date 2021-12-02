package middlewares

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

func NewEndpointBackgroundCronLoggerMiddleware(logger log.Logger, endpointMethod string) endpoint.Middleware {
	return MakeEndpointBackgroundCronLoggerMiddleware(logger, endpointMethod)
}