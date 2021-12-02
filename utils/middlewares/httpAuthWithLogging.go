package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/mytrix-technology/mylibgo/utils/helper"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/ua-parser/uap-go/uaparser"
)

type AuthWithLoggingMiddlewareFunc func(next http.Handler, apiCode string, menuCode string) http.Handler

// MakeUTHttpAuthMiddleware create MiddlewareFunc for go http.Handler
func MakeHttpTransportLoggingWithAuthMiddleware(idmAuthorizeUrl string, netClient *http.Client, utKey string, logger log.Logger) AuthWithLoggingMiddlewareFunc {
	return func(next http.Handler, apiCode string, menuCode string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqid := r.Header.Get("X-Request-Id")
			traceid := r.Header.Get("X-Trace-Id")

			ua := uaparser.NewFromSaved()
			cl := ua.Parse(r.Header.Get("User-Agent"))

			//if r.Method == "OPTIONS" {
			//	w.WriteHeader(200)
			//	return
			//}

			lvl := "INFO"
			message := ""
			allowed := true
			//headerJson, _ := json.Marshal(r.Header)
			contentType := r.Header.Get(helper.HeaderContentType)
			debugger := level.Debug(logger).Log

			defer func() {
				_ = logger.Log(
					"request-id", reqid,
					"trace-id", traceid,
					"lvl", lvl,
					"event", "incoming request",
					"uri", r.RequestURI,
					"method", r.Method,
					"headers", r.Header,
					"msg", message,
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

			// get token
			var auth = r.Header.Get("Authorization")
			if auth == "" {
				lvl = "ERROR"
				message = "forbidden: Authorization header not found"
				allowed = false
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// validate token
			var tokens = strings.Split(auth, "=")
			if len(tokens) == 1 || len(tokens) == 0 {
				lvl = "ERROR"
				message = "forbidden: invalid token"
				allowed = false
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if strings.TrimSpace(tokens[0]) != "token" || strings.TrimSpace(tokens[1]) == "" {
				lvl = "ERROR"
				message = "forbidden: invalid token"
				allowed = false
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			var token = strings.TrimSpace(tokens[1])

			var datetime = ""
			var signature = ""
			var bodyBytes []byte
			var bodyRead = false

			if r.Method == "POST" && contentType == helper.HttpContentTypeJson {
				if r.Body != nil {
					bodyBytes, _ = ioutil.ReadAll(r.Body)
				}
				bodyRead = true

				var authBody AuthBody
				err := json.Unmarshal(bodyBytes, &authBody)
				//if err != nil {
				//	lvl = "ERROR"
				//	message = fmt.Sprintf("forbidden: invalid datetime & signature. %s", err)
				//	allowed = false
				//	http.Error(w, "Forbidden", http.StatusForbidden)
				//	return
				//}

				if err == nil {
					datetime = authBody.Datetime
					if datetime != "" {
						r.Header.Set("datetime", datetime)
					}

					signature = authBody.Signature
					if signature != "" {
						r.Header.Set("signature", signature)
					}
				}
			}

			if datetime == "" {
				datetime = r.Header.Get("datetime")
			}

			if signature == "" {
				signature = r.Header.Get("signature")
			}

			if datetime == "" {
				lvl = "ERROR"
				message = "forbidden: datetime header not found"
				allowed = false
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if signature == "" {
				if datetime == "" {
					lvl = "ERROR"
					message = "forbidden: signature header not found"
					allowed = false
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}

			// validate signature
			if !validate(utKey, token, datetime, signature) {
				lvl = "ERROR"
				message = "forbidden: invalid signature"
				allowed = false
				http.Error(w, "Invalid Signature", http.StatusForbidden)
				return
			}

			authReq := RequestBody{
				RoleCode:  "",
				MenuCode:  "",
				APICode:   apiCode,
				Date:      datetime,
				Signature: signature,
				Data:      string(bodyBytes),
			}
			// Restore the io.ReadCloser to its original state if already read
			if bodyRead {
				r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(authReq); err != nil {
				lvl = "ERROR"
				message = "internal error: failed to serialize internal request"
				allowed = false
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			request, err := http.NewRequest("POST", idmAuthorizeUrl, &buf)
			if err != nil {
				lvl = "ERROR"
				message = "internal error: failed to create http request object"
				allowed = false
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			request.Header = r.Header
			request.Header.Set("Content-Type", "application/json")

			_ = debugger("event", "auth call", "msg", "sending auth request to IDM", "payload", authReq)
			resp, err := netClient.Do(request)
			if err != nil {
				lvl = "ERROR"
				message = fmt.Sprintf("internal error: failed to contact IDM: %s", err)
				allowed = false
				http.Error(w, fmt.Sprintf("failed to contact IDM: %s", err), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == 403 {
				lvl = "INFO"
				message = "forbidden"
				allowed = false
				w.WriteHeader(http.StatusForbidden)
				res := map[string]interface{}{
					"code":  http.StatusForbidden,
					"error": "forbidden",
					"data":  nil,
				}

				resJson, err := json.Marshal(res)
				if err != nil {
					return
				}

				w.Write(resJson)
				return
			}

			if resp.StatusCode > 400 {
				lvl = "ERROR"
				message = fmt.Sprintf("failed to authenticate request to IDM: %v", resp)
				allowed = false
				w.WriteHeader(http.StatusInternalServerError)
				res := map[string]interface{}{
					"code":  http.StatusInternalServerError,
					"error": "failed to authenticate request to IDM",
					"data":  nil,
				}

				resJson, err := json.Marshal(res)
				if err != nil {
					return
				}

				w.Write(resJson)
				return
			}

			var authResp AuthorizeResponse
			if err := decodeJSONBody(&authResp, resp.Body); err != nil {
				lvl = "ERROR"
				message = "internal error: failed to decode IDM response"
				allowed = false
				http.Error(w, "failed to decode IDM response", http.StatusInternalServerError)
				return
			}

			if !authResp.Success {
				lvl = "ERROR"
				message = fmt.Sprintf("unauthorized: IDM Error: %s", authResp.Message)
				allowed = false
				w.WriteHeader(http.StatusForbidden)
				res := map[string]interface{}{
					"code":  http.StatusForbidden,
					"error": fmt.Sprintf("failed to authenticate request to IDM: %s", authResp.Message),
					"data":  nil,
				}

				resJson, err := json.Marshal(res)
				if err != nil {
					return
				}

				w.Write(resJson)
				return
			}

			if !authResp.Data.Authorize {
				lvl = "ERROR"
				message = fmt.Sprintf("unauthorized: user not authorized: %s", err)
				allowed = false
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			userID := strconv.FormatInt(authResp.Data.UserID, 10)

			// set userID
			r.Header.Set("UT-UserID", userID)
			r.Header.Set("UT-Token", token)

			ctx := context.WithValue(r.Context(), "UT-UserID", authResp.Data.UserID)
			ctx = context.WithValue(ctx, helper.ContextKeyUTPrivilegeElement, authResp.Data.PrivilegeElement)
			r = r.Clone(ctx)
			next.ServeHTTP(w, r)

			r.Header.Del("UT-UserID")
			r.Header.Del("UT-Token")
		})
	}
}

func MakeNoopHttpTransportLoggingWithAuthMiddleware(idmAuthorizeUrl string, netClient *http.Client, utKey string, logger log.Logger) AuthWithLoggingMiddlewareFunc {
	return func(next http.Handler, apiCode string, menuCode string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqid := r.Header.Get("X-Request-Id")
			traceid := r.Header.Get("X-Trace-Id")

			ua := uaparser.NewFromSaved()
			cl := ua.Parse(r.Header.Get("User-Agent"))

			level := "INFO"
			message := ""
			allowed := true
			//headerJson, _ := json.Marshal(r.Header)

			defer func() {
				_ = logger.Log(
					"request-id", reqid,
					"trace-id", traceid,
					"level", level,
					"event", "incoming request",
					"uri", r.RequestURI,
					"method", r.Method,
					"headers", r.Header,
					"msg", message,
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

			// get token
			var auth = r.Header.Get("Authorization")
			if auth == "" {
				level = "ERROR"
				message = "forbidden: Authorization header not found"
				allowed = false
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// validate token
			var tokens = strings.Split(auth, "=")
			if len(tokens) == 1 || len(tokens) == 0 {
				level = "ERROR"
				message = "forbidden: invalid token"
				allowed = false
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if strings.TrimSpace(tokens[0]) != "token" || strings.TrimSpace(tokens[1]) == "" {
				level = "ERROR"
				message = "forbidden: invalid token"
				allowed = false
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			var token = strings.TrimSpace(tokens[1])

			var datetime = ""
			var signature = ""
			var bodyBytes []byte

			if r.Method == "POST" {
				if r.Body != nil {
					bodyBytes, _ = ioutil.ReadAll(r.Body)
				}

				var authBody AuthBody
				err := json.Unmarshal(bodyBytes, &authBody)
				//if err != nil {
				//	level = "ERROR"
				//	message = fmt.Sprintf("forbidden: invalid datetime & signature. %s", err)
				//	allowed = false
				//	http.Error(w, "Forbidden", http.StatusForbidden)
				//	return
				//}

				if err == nil {
					datetime = authBody.Datetime
					if datetime != "" {
						r.Header.Set("datetime", datetime)
					}

					signature = authBody.Signature
					if signature != "" {
						r.Header.Set("signature", signature)
					}
				}
			}

			if datetime == "" {
				datetime = r.Header.Get("datetime")
			}

			if signature == "" {
				signature = r.Header.Get("signature")
			}

			if datetime == "" {
				level = "ERROR"
				message = "forbidden: datetime header not found"
				allowed = false
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if signature == "" {
				if datetime == "" {
					level = "ERROR"
					message = "forbidden: signature header not found"
					allowed = false
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}

			// validate signature
			if !validate(utKey, token, datetime, signature) {
				level = "ERROR"
				message = "forbidden: invalid signature"
				allowed = false
				http.Error(w, "Invalid Signature", http.StatusForbidden)
				return
			}

			authReq := RequestBody{
				RoleCode:  "",
				MenuCode:  "",
				APICode:   apiCode,
				Date:      datetime,
				Signature: signature,
				Data:      string(bodyBytes),
			}
			// Restore the io.ReadCloser to its original state
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(authReq); err != nil {
				level = "ERROR"
				message = "internal error: failed to serialize internal request"
				allowed = false
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			request, err := http.NewRequest("POST", idmAuthorizeUrl, &buf)
			if err != nil {
				level = "ERROR"
				message = "internal error: failed to create http request object"
				allowed = false
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			request.Header = r.Header

			// set userID
			r.Header.Set("UT-UserID", "999999")
			r.Header.Set("UT-Token", token)

			ctx := context.WithValue(r.Context(), "UT-UserID", 999999)
			r = r.Clone(ctx)
			next.ServeHTTP(w, r)

			r.Header.Del("UT-UserID")
			r.Header.Del("UT-Token")
		})
	}
}
