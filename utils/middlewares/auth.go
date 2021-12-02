package middlewares

import (
	"net/http"

	"github.com/go-kit/kit/log"

	"github.com/mytrix-technology/mylibgo/security/cryptography"
)

type AuthorizeRequest struct {
	Token    string `json:"token"`
	RoleCode string `json:"role_code"`
	MenuCode string `json:"menu_code"`
	APICode  string `json:"api_code"`
}

type AuthorizeResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    *AuthorizeData `json:"data"`
}

type AuthorizeData struct {
	UserID    			int64  		`json:"user_id"`
	UserCode  			string 		`json:"user_code"`
	Token     			string 		`json:"token"`
	Authorize 			bool   		`json:"authorize"`
	PrivilegeElement	[]string	`json:"privilege_element"`
}

type RequestBody struct {
	RoleCode  string `json:"role_code"`
	MenuCode  string `json:"menu_code"`
	APICode   string `json:"api_code"`
	Date      string `json:"datetime"`
	Signature string `json:"signature"`
	Data 	  string `json:"data"`
}

func NewEndpointAuthMiddlewareFunc(idmAuthorizeUrl string, netClient *http.Client, utKey string) AuthEndpointMiddlewareFunc {
	return MakeEndpointAuthMiddlewareFunc(idmAuthorizeUrl, netClient, utKey)
}

func NewHttpLoggingWithAuthMiddlewareFunc(idmAuthorizeUrl string, netClient *http.Client, utKey string, logger log.Logger) AuthWithLoggingMiddlewareFunc {
	return MakeHttpTransportLoggingWithAuthMiddleware(idmAuthorizeUrl, netClient, utKey, logger)
}

func NewHttpLoggingWithAuthMiddlewarePublicFunc(logger log.Logger) AuthWithLoggingMiddlewarePublicFunc {
	return MakeHttpTransportLoggingWithAuthMiddlewarePublic(logger)
}

func validate(utkey, token, datetime, signature string) bool {
	encoded := cryptography.EncodeSHA256HMAC(utkey, token, datetime)
	if encoded == signature {
		return true
	}
	return false
}
