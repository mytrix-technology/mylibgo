package middlewares

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/mytrix-technology/mylibgo/utils/helper"
)

type CORSOption func(*CorsConfig)

type CorsConfig struct {
	// AllowOrigin defines a list of origins that may access the resource.
	// Optional. Default value []string{"*"}.
	AllowOrigins []string `yaml:"allow_origins"`

	// AllowMethods defines a list methods allowed when accessing the resource.
	// This is used in response to a preflight request.
	// Optional. Default value DefaultCORSConfig.AllowMethods.
	AllowMethods []string `yaml:"allow_methods"`

	// AllowHeaders defines a list of request headers that can be used when
	// making the actual request. This is in response to a preflight request.
	// Optional. Default value []string{}.
	AllowHeaders []string `yaml:"allow_headers"`

	// AllowCredentials indicates whether or not the response to the request
	// can be exposed when the credentials flag is true. When used as part of
	// a response to a preflight request, this indicates whether or not the
	// actual request can be made using credentials.
	// Optional. Default value false.
	AllowCredentials bool `yaml:"allow_credentials"`

	// ExposeHeaders defines a whitelist headers that clients are allowed to
	// access.
	// Optional. Default value []string{}.
	ExposeHeaders []string `yaml:"expose_headers"`

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached.
	// Optional. Default value 0.
	MaxAge int `yaml:"max_age"`
}

var (
	// DefaultCORSConfig is the default CORS middlewares config.
	DefaultCORSConfig = CorsConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
)

func CorsAllowOrigin(origins []string) CORSOption {
	return func(config *CorsConfig) {
		if len(origins) > 0 {
			config.AllowOrigins = origins
		}
	}
}

func CorsAllowMethods(methods []string) CORSOption {
	return func(config *CorsConfig) {
		if len(methods) > 0 {
			config.AllowMethods = methods
		}
	}
}

func CorsAllowHeaders(headers []string) CORSOption {
	return func(config *CorsConfig) {
		if len(headers) > 0 {
			config.AllowHeaders = headers
		}
	}
}

func CorsAllowCredentials(allowCredential bool) CORSOption {
	return func(config *CorsConfig) {
		config.AllowCredentials = allowCredential
	}
}

func CorsMaxAge(maxAge int) CORSOption {
	return func(config *CorsConfig) {
		config.MaxAge = maxAge
	}
}

func NewCORSConfig(options ...CORSOption) CorsConfig {
	c := CorsConfig{}
	for _, op := range options {
		op(&c)
	}

	return c
}

func CORSMiddleware(options ...CORSOption) MiddlewareFunc {
	config := DefaultCORSConfig
	for _, op := range options {
		op(&config)
	}

	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")
	maxAge := strconv.Itoa(config.MaxAge)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get(helper.HeaderOrigin)
			allowOrigin := ""

			// Check allowed origins
			for _, o := range config.AllowOrigins {
				if o == "*" && config.AllowCredentials {
					allowOrigin = origin
					break
				}
				if o == "*" || o == origin {
					allowOrigin = o
					break
				}
				if matchSubdomain(origin, o) {
					allowOrigin = origin
					break
				}
			}

			log.Println("[cors] origin: ", allowOrigin)

			if r.Method != http.MethodOptions {
				w.Header().Add(helper.HeaderVary, helper.HeaderOrigin)
				w.Header().Set(helper.HeaderAccessControlAllowOrigin, allowOrigin)
				if config.AllowCredentials {
					w.Header().Set(helper.HeaderAccessControlAllowCredentials, "true")
				}
				if exposeHeaders != "" {
					w.Header().Set(helper.HeaderAccessControlExposeHeaders, exposeHeaders)
				}

				next.ServeHTTP(w, r)
				return
			}

			// Handling pre-flight request
			log.Println("[cors] handling preflight")
			w.Header().Add(helper.HeaderVary, helper.HeaderOrigin)
			w.Header().Add(helper.HeaderVary, helper.HeaderAccessControlRequestMethod)
			w.Header().Add(helper.HeaderVary, helper.HeaderAccessControlRequestHeaders)
			w.Header().Set(helper.HeaderAccessControlAllowOrigin, allowOrigin)
			w.Header().Set(helper.HeaderAccessControlAllowMethods, allowMethods)

			if config.AllowCredentials {
				w.Header().Set(helper.HeaderAccessControlAllowCredentials, "true")
			}

			if allowHeaders != "" {
				w.Header().Set(helper.HeaderAccessControlAllowHeaders, allowHeaders)
			} else {
				h := r.Header.Get(helper.HeaderAccessControlRequestHeaders)
				if h != "" {
					w.Header().Set(helper.HeaderAccessControlAllowHeaders, h)
				}
			}
			if config.MaxAge > 0 {
				w.Header().Set(helper.HeaderAccessControlMaxAge, maxAge)
			}

			log.Printf("[cors] written header: %+v\n", w.Header())
			w.WriteHeader(http.StatusOK)
		})
	}
}
