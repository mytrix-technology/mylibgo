package middlewares

import (
	"time"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)
func MakeEndpointBackgroundCronLoggerMiddleware(logger log.Logger, cronName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (result interface{}, err error){
			logger = log.With(logger, "method", "cron", "cron-name", cronName)
			_ = level.Info(logger).Log("event", "job triggered", "msg", fmt.Sprintf("starting job for %s", cronName))
			go func() {
				defer func(begin time.Time){
					logger = log.With(logger, "event", "job ended", "duration", time.Since(begin))
					if err != nil {
						_ = level.Error(logger).Log("msg", err.Error())
					} else {
						_ = level.Info(logger).Log("msg", "success")
					}
				}(time.Now())
				result, err = next(ctx, request)
			}()
			return nil, nil
		}
	}
}
