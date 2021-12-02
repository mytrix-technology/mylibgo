package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/mytrix-technology/mylibgo/messaging/slack"
)

func MakePushSlackErrNotiicationfMiddleware(webhookUrl, host, version string, isSendNotif bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if isSendNotif {
					err := recovery(recover())
					if err != nil {
						s := slack.Slack{
							Host:       host,
							Version:    version,
							WebhookURL: webhookUrl,
							Message:    err.Error(),
							StackTrace: debug.Stack(),
						}

						s.PostError()
					}
				}

			}()

			next.ServeHTTP(w, r)
		})
	}
}
