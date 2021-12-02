package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/mytrix-technology/mylibgo/networking/http/helper"
	"github.com/mytrix-technology/mylibgo/utils/types/errs"
)

type AuthEndpointMiddlewareFunc func(apiCode string, menuCode string) endpoint.Middleware

func MakeEndpointAuthMiddlewareFunc(idmAuthorizeUrl string, netClient *http.Client, utKey string) AuthEndpointMiddlewareFunc {
	return func(apiCode string, menuCode string) endpoint.Middleware {
		return func(next endpoint.Endpoint) endpoint.Endpoint {
			return func(ctx context.Context, request interface{}) (result interface{}, err error) {
				result = nil
				////get token
				//auth, ok := ctx.Value(helper.ContextKeyRequestAuthorization).(string)
				//if !ok {
				//	err = fmt.Errorf("%w: authorization context not found", errs.ErrForbidden)
				//	return
				//}
				//
				//if auth == "" {
				//	err = fmt.Errorf("%w: authorization context is empty", errs.ErrForbidden)
				//	return
				//}
				//
				//// validate token
				//var tokens = strings.Split(auth, "=")
				//if len(tokens) == 1 || len(tokens) == 0 {
				//	err = fmt.Errorf("%w: token error", errs.ErrForbidden)
				//	return
				//}
				//
				//if strings.TrimSpace(tokens[0]) != "token" || strings.TrimSpace(tokens[1]) == "" {
				//	err = fmt.Errorf("%w: token error", errs.ErrForbidden)
				//	return
				//}
				//
				//var token = strings.TrimSpace(tokens[1])

				//// get datetime
				//datetime, ok := ctx.Value(helper.ContextKeyRequestDatetime).(string)
				//if !ok {
				//	err = fmt.Errorf("%w: datetime context not found", errs.ErrForbidden)
				//	return
				//}
				//
				//if datetime == "" {
				//	err = fmt.Errorf("%w: datetime context is empty", errs.ErrForbidden)
				//	return
				//}
				//
				//// get signature
				//signature, ok := ctx.Value(helper.ContextKeyRequestSignature).(string)
				//if !ok {
				//	err = fmt.Errorf("%w: signature context not found", errs.ErrForbidden)
				//	return
				//}
				//
				//if signature == "" {
				//	err = fmt.Errorf("%w: signature context is empty", errs.ErrForbidden)
				//	return
				//}

				var auth *helper.AuthObject
				auth, err = helper.AuthObjectFromContext(ctx)
				if err != nil {
					return
				}

				// validate signature
				if !validate(utKey, auth.Token, auth.Datetime, auth.Signature) {
					err = fmt.Errorf("%w: invalid signature", errs.ErrForbidden)
					return
				}

				authReq := RequestBody{
					RoleCode:  "",
					MenuCode:  "",
					APICode:   apiCode,
					Date:      auth.Datetime,
					Signature: auth.Signature,
				}

				var buf *bytes.Buffer
				if err = json.NewEncoder(buf).Encode(authReq); err != nil {
					err = fmt.Errorf("%s. %w: token error", errs.ErrInternalServerError, errs.ErrForbidden)
					return
				}

				req, localErr := http.NewRequest("POST", idmAuthorizeUrl, buf)
				if localErr != nil {
					err = fmt.Errorf("%s. %w: authorization error", errs.ErrInternalServerError, errs.ErrForbidden)
					return
				}

				req.Header.Set("Authorization", "token="+auth.Token)
				req.Header.Set("datetime", auth.Datetime)
				req.Header.Set("signature", auth.Signature)

				resp, localErr := netClient.Do(req)
				if localErr != nil {
					err = fmt.Errorf("%s. %w: failed to contact IDM: %s", errs.ErrInternalServerError, errs.ErrForbidden, err)
					return
				}
				defer resp.Body.Close()

				var authResp AuthorizeResponse
				if err = json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
					err = fmt.Errorf("%s. %w: failed to decode IDM response", errs.ErrInternalServerError, errs.ErrForbidden)
					return
				}

				// set context
				ctx = context.WithValue(ctx, "UT-UserID", authResp.Data.UserID)
				ctx = context.WithValue(ctx, "UT-Token", auth.Token)

				result, err = next(ctx, request)
				return
			}
		}
	}
}
