package middlewares

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/mytrix-technology/mylibgo/utils/helper"
	"github.com/mytrix-technology/mylibgo/messaging/slack"
)

type PanicNotifPoster interface{
	Post(err error, stackTrace string) error
}


type slackPanicNotifPoster struct {
	si *slack.SlackIntegration
}

func MakeSlackNotifPoster(si *slack.SlackIntegration) PanicNotifPoster {
	return &slackPanicNotifPoster{si}
}

func (sn *slackPanicNotifPoster) Post(err error, stackTrace string) error {
	msg := err.Error()
	stackFormatted := "```" + stackTrace + "```"
	return sn.si.PostError(msg, slack.AddTags("panic"), slack.AddAttachment(stackFormatted, slack.WithTitle("Stack Trace")))
}

func MakeErrPanicNotificationMiddleware(poster PanicNotifPoster, debugLogger helper.DebugFieldLogger) func(http.Handler) http.Handler {
	if debugLogger == nil {
		debugLogger = helper.CreateNoopFieldLogger()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recovery(recover())
				if err != nil {
					stackTrace := string(debug.Stack())
					if r := poster.Post(err, stackTrace); r != nil {
						_ = debugLogger("event", "Panic Notification failed", "error", r.Error())
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func recovery(r interface{}) error {
	var err error
	if r != nil {
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("unknown error")
		}
	}
	return err
}

