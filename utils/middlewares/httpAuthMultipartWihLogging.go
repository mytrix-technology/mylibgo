package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/ua-parser/uap-go/uaparser"
)

type AuthMultipartWithLoggingMiddlewareFunc func(next http.Handler, apiCode string, menuCode string) http.Handler


func MakeHttpTransportLoggingWithAuthMultipartMiddleware(idmAuthorizeUrl string, netClient *http.Client, utKey string, logger log.Logger) AuthWithLoggingMiddlewareFunc {
	return func(next http.Handler, apiCode string, menuCode string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			reqid := r.Header.Get("X-Request-Id")
			traceid := r.Header.Get("X-Trace-Id")

			ua := uaparser.NewFromSaved()
			cl := ua.Parse(r.Header.Get("User-Agent"))

			//if r.Method == "OPTIONS" {
			//	w.WriteHeader(200)
			//	return
			//}

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
			if r.Method == "POST" {
				if r.Form == nil {
					r.ParseMultipartForm(32 << 20)
				}

				datetime = r.FormValue("datetime")
				if datetime != "" {
					r.Header.Set("datetime", datetime)
				}

				signature = r.FormValue("signature")
				if signature != "" {
					r.Header.Set("signature", signature)
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
				//Data:  	   string(bodyBytes),
			}

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

			buffHeader := r.Header.Clone()
			request.Header = http.Header{}
			request.Header.Set("Authorization", r.Header.Get("Authorization"))
			request.Header.Set("datetime", r.Header.Get("datetime"))
			request.Header.Set("signature", r.Header.Get("signature"))

			resp, err := netClient.Do(request)
			if err != nil {
				level = "ERROR"
				message = fmt.Sprintf("internal error: failed to contact IDM: %s", err)
				allowed = false
				http.Error(w, fmt.Sprintf("failed to contact IDM: %s", err), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			// Restore the io.ReadCloser to its original state
			r.Header = buffHeader

			if resp.StatusCode == 403 {
				level = "INFO"
				message = "forbidden"
				allowed = false
				w.WriteHeader(http.StatusForbidden)
				res := map[string]interface{}{
					"code": http.StatusForbidden,
					"error": "forbidden",
					"data": nil,
				}

				resJson, err := json.Marshal(res)
				if err != nil {
					return
				}

				w.Write(resJson)
				return
			}

			if resp.StatusCode > 400 {
				level = "ERROR"
				message = fmt.Sprintf("failed to authenticate request to IDM: %v", resp)
				allowed = false
				w.WriteHeader(http.StatusInternalServerError)
				res := map[string]interface{}{
					"code": http.StatusInternalServerError,
					"error": "failed to authenticate request to IDM",
					"data": nil,
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
				level = "ERROR"
				message = "internal error: failed to decode IDM response"
				allowed = false
				http.Error(w, fmt.Sprintf("failed to decode IDM response"), http.StatusInternalServerError)
				return
			}

			if !authResp.Success {
				level = "ERROR"
				message = fmt.Sprintf("unauthorized: IDM Error: %s", authResp.Message)
				allowed = false
				w.WriteHeader(http.StatusForbidden)
				res := map[string]interface{}{
					"code": http.StatusForbidden,
					"error": fmt.Sprintf("failed to authenticate request to IDM: %s", authResp.Message),
					"data": nil,
				}

				resJson, err := json.Marshal(res)
				if err != nil {
					return
				}

				w.Write(resJson)
				return
			}

			if !authResp.Data.Authorize {
				level = "ERROR"
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
			r = r.Clone(ctx)
			next.ServeHTTP(w, r)

			r.Header.Del("UT-UserID")
			r.Header.Del("UT-Token")
		})
	}
}