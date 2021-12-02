package middlewares

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/rs/xid"
	"github.com/ua-parser/uap-go/uaparser"
)

func HttpRequestIDInjectorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqid := xid.New().String()
		traceid := r.Header.Get("X-Trace-Id")
		if traceid == "" {
			traceid = xid.New().String()
		}

		r.Header.Set("X-Request-Id", reqid)
		r.Header.Set("X-Trace-Id", traceid)

		// inject header
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Origin,Accept,datetime,signature")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func MakeHttpRequestBodySizeLimiterMiddleware(limitSizeInByte int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, limitSizeInByte)
			next.ServeHTTP(w, r)
		})
	}
}

func MakeHttpTransportLoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqid := r.Header.Get("X-Request-Id")
			traceid := r.Header.Get("X-Trace-Id")

			ua := uaparser.NewFromSaved()
			cl := ua.Parse(r.Header.Get("User-Agent"))
			flogger := log.With(logger, "request-id", reqid, "trace-id", traceid)
			_ = level.Info(flogger).Log(
				"event", "incoming request",
				"uri", r.RequestURI,
				"method", r.Method,
				"headers", r.Header,
				"origin", r.Header.Get("X-Forwarded-For"),
				"protocol", r.Proto,
				"user-agent", cl.UserAgent.ToString(),
				"device", cl.Device.ToString(),
				"os", cl.Os.ToString(),
			)
			next.ServeHTTP(w, r)
		})
	}
}
