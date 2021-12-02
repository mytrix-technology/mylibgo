package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/ua-parser/uap-go/uaparser"

	"github.com/mytrix-technology/mylibgo/security/cryptography"
	"github.com/mytrix-technology/mylibgo/utils/datetime"
)

const (
	PubUTSignKey = "4utrwFjnQAF9yU6P"
)

type AuthWithLoggingMiddlewarePublicFunc func(next http.Handler) http.Handler

func MakeHttpTransportLoggingWithAuthMiddlewarePublic(logger log.Logger) AuthWithLoggingMiddlewarePublicFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqid := r.Header.Get("X-Request-Id")
			traceid := r.Header.Get("X-Trace-Id")

			ua := uaparser.NewFromSaved()
			cl := ua.Parse(r.Header.Get("User-Agent"))

			level := "INFO"
			messagelog := ""
			allowed := true

			defer func() {
				_ = logger.Log(
					"request-id", reqid,
					"trace-id", traceid,
					"level", level,
					"event", "incoming request",
					"uri", r.RequestURI,
					"method", r.Method,
					"headers", r.Header,
					"msg", messagelog,
					"auth-status", allowed,
					"content-type", r.Header.Get("Content-Type"),
					"content-length", r.Header.Get("Content-Length"),
					"origin", r.RemoteAddr,
					"protocol", r.Proto,
					"user-agent", cl.UserAgent.ToString(),
					"device", cl.Device.ToString(),
					"os", cl.Os.ToString(),
				)
			}()

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization,Origin,Accept,API-Key,Source,DeviceID,datetime,signature")
			w.Header().Set("Content-Type", "application/json")

			if r.Method == "OPTIONS" {
				w.WriteHeader(200)
				return
			}

			// get token
			var auth = r.Header.Get("Authorization")
			if auth == "" {
				level = "ERROR"
				messagelog = "forbidden: Authorization header not found"
				allowed = false
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// validate token
			var tokens = strings.SplitN(auth, "=", 2)
			if len(tokens) == 1 || len(tokens) == 0 {
				level = "ERROR"
				messagelog = "forbidden: invalid token"
				allowed = false
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if strings.TrimSpace(tokens[0]) != "token" || strings.TrimSpace(tokens[1]) == "" {
				level = "ERROR"
				messagelog = "forbidden: invalid token"
				allowed = false
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			var token string

			// some utilities to generate token
			if len(tokens) >= 2 {
				for index := 1; index < len(tokens); index++ {
					token = token + tokens[index]
				}
			}
			token = strings.TrimSpace(token)
			if token == "" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// decode token base64 = apiKey:deviceID
			decodeToken, err := cryptography.DecodeBASE64(token)
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			var clientAuths = strings.Split(string(decodeToken), ":")
			if len(clientAuths) == 1 || len(clientAuths) == 0 {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// get apiKeyReq
			var apiKeyReq = r.Header.Get("API-Key")
			if apiKeyReq == "" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// get sourceReq
			var sourceReq = r.Header.Get("Source")
			if sourceReq == "" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// get deviceIDReq
			var deviceIDReq = r.Header.Get("DeviceID")
			if deviceIDReq == "" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// get datetimeReq
			var datetimeReq = r.Header.Get("datetime")
			if datetimeReq == "" || datetime.TimeToStringPattern(datetime.StringToTime(datetimeReq), "2006-01-02") != datetime.TimeToStringPattern(time.Now(), "2006-01-02") {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			var apiKey = strings.TrimSpace(clientAuths[0])
			var deviceID = strings.TrimSpace(clientAuths[1])
			var message = cryptography.EncodeBASE64([]byte(apiKeyReq + ":" + sourceReq + ":" + deviceIDReq + ":" + datetimeReq))

			if apiKey != apiKeyReq {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if deviceID != deviceIDReq {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if r.Method == "GET" {
				// get signature
				var signatureReq = r.Header.Get("signature")
				if signatureReq == "" {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}

				// validate signature
				var signature = cryptography.EncodeSHA256HMAC(PubUTSignKey, apiKey, deviceID, message, datetimeReq)
				if signature != signatureReq {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}

			r.Header.Set("UT-APIKey", apiKeyReq)
			r.Header.Set("UT-DeviceID", deviceIDReq)
			r.Header.Set("UT-Message", message)

			ctx := context.WithValue(r.Context(), "UserContextKey", nil)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

			r.Header.Del("UT-APIKey")
			r.Header.Del("UT-DeviceID")
			r.Header.Del("UT-Message")
		})
	}
}
