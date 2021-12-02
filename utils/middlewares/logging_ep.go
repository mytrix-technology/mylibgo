package middlewares

import (
	"context"
	"time"

	"github.com/mytrix-technology/mylibgo/networking/http/helper"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type AppendKeyvalser interface {
	AppendKeyvals(keyvals []interface{}) []interface{}
}

const (
	tookKey = "duration"
)

func MakeEndpointLoggingMiddleware(logger log.Logger, endPointMethod string) endpoint.Middleware {
	if logger == nil {
		return nil
	}
	epLogger := log.With(logger, "method", endPointMethod)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (result interface{}, err error) {
			reqid, ok := helper.ReqIDFromContext(ctx)
			if !ok {
				reqid = ""
			}

			traceid, ok := helper.TraceIDFromContext(ctx)
			if !ok {
				reqid = ""
			}

			flogger := log.With(epLogger, "request-id", reqid, "trace-id", traceid)
			// _ = level.Info(flogger).Log("event", "request processing", "input", request)
			defer func(begin time.Time) {
				var keyVals []interface{}
				if err != nil {
					keyVals = makeKeyvals(request, result, "event", "error", "msg", err.Error(), tookKey, time.Since(begin))
					_ = level.Error(flogger).Log(keyVals...)
				} else {
					keyVals = makeKeyvals(request, result)
					_ = level.Info(flogger).Log("event", "success", "duration", time.Since(begin))
					if keyVals != nil {
						_ = level.Debug(flogger).Log(keyVals...)
					}
				}
			}(time.Now())

			result, err = next(ctx, request)
			return
		}
	}
}

func makeKeyvals(req, resp interface{}, keyVals ...interface{}) []interface{} {
	KVs := keyVals
	notNil := false
	if l, ok := req.(AppendKeyvalser); ok {
		notNil = true
		KVs = l.AppendKeyvals(KVs)
	}
	if l, ok := resp.(AppendKeyvalser); ok {
		notNil = true
		KVs = l.AppendKeyvals(KVs)
	}

	if notNil {
		return KVs
	}

	return nil
}
