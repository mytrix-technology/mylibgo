package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/mytrix-technology/mylibgo/utils/helper"
)

type AuthMiddlewareFunc func(next http.Handler, apiCode string, menuCode string) http.Handler

type AuthBody struct {
	Datetime  string `json:"datetime"`
	Signature string `json:"signature"`
}

// MakeUTHttpAuthMiddleware create MiddlewareFunc for go http.Handler
func MakeHttpTransportAuthMiddleware(idmAuthorizeUrl string, netClient *http.Client, utKey string) AuthMiddlewareFunc {
	return func(next http.Handler, apiCode string, menuCode string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// get token
			var auth = r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// validate token
			var tokens = strings.Split(auth, "=")
			if len(tokens) == 1 || len(tokens) == 0 {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if strings.TrimSpace(tokens[0]) != "token" || strings.TrimSpace(tokens[1]) == "" {
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
				if err != nil {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}

				datetime = authBody.Datetime
				if datetime != "" {
					r.Header.Set("datetime", datetime)
				}

				signature = authBody.Signature
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
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if signature == "" {
				if datetime == "" {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}

			// validate signature
			if !validate(utKey, token, datetime, signature) {
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			request, err := http.NewRequest("POST", idmAuthorizeUrl, &buf)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			request.Header = r.Header

			resp, err := netClient.Do(request)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to contact IDM: %s", err), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == 403 {
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
				http.Error(w, fmt.Sprintf("failed to decode IDM response"), http.StatusInternalServerError)
				return
			}

			if !authResp.Success {
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

func decodeJSONBody(dest interface{}, body io.Reader) error {
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	if err := json.Unmarshal(buf, dest); err != nil {
		return fmt.Errorf("failed to parse response body: %s", err)
	}

	return nil
}
