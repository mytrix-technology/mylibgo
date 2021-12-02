package router

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/mytrix-technology/mylibgo/utils/helper"
	"github.com/mytrix-technology/mylibgo/utils/middlewares"
	"github.com/julienschmidt/httprouter"
)

//type MiddlewareFunc func(next http.Handler) http.Handler

type Router struct {
	router      *httprouter.Router
	isInit      bool
	prefix      string
	middlewares []middlewares.MiddlewareFunc
	routes      []*Route
	subRouters  []*Router
	cors        *corsConfig
	debugLogger helper.DebugFieldLogger
}

type Route struct {
	router  *Router
	path    string
	methods []string
	handler http.Handler
}

type corsConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     string   `yaml:"allow_methods"`
	AllowHeaders     string   `yaml:"allow_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`

	// ExposeHeaders defines a whitelist headers that clients are allowed to
	// access.
	// Optional. Default value []string{}.
	ExposeHeaders string `yaml:"expose_headers"`
	MaxAge        int    `yaml:"max_age"`
	MaxAgeStr     string
}

var (
	// DefaultCORSConfig is the default CORS middlewares config.
	DefaultCORSConfig = corsConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: strings.Join([]string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}, ","),
	}
)

type CorsOption func(*corsConfig)
type RouterOption func(*Router)

func NewRouter(options ...RouterOption) *Router {
	router := &Router{
		router: httprouter.New(),
		//methods: []string{},
		middlewares: []middlewares.MiddlewareFunc{},
		routes:      []*Route{},
		subRouters:  []*Router{},
		isInit:      false,
		debugLogger: helper.CreateNoopFieldLogger(),
	}

	for _, op := range options {
		op(router)
	}
	return router
}

func SetCORSConfig(options ...CorsOption) RouterOption {
	config := DefaultCORSConfig
	return func(r *Router) {
		for _, op := range options {
			op(&config)
		}
		r.cors = &config
	}
}

func SetDebugLogger(debugLogger helper.DebugFieldLogger) RouterOption {
	return func(r *Router) {
		if debugLogger != nil {
			r.debugLogger = debugLogger
		}
	}
}

func (rtr *Router) Use(middlewares ...middlewares.MiddlewareFunc) {
	rtr.middlewares = append(rtr.middlewares, middlewares...)
}

func (rtr *Router) Subroute(pathPrefix string) *Router {
	router := Router{
		prefix:      pathPrefix,
		router:      rtr.router,
		middlewares: []middlewares.MiddlewareFunc{},
		routes:      []*Route{},
		subRouters:  []*Router{},
		isInit:      false,
	}

	rtr.subRouters = append(rtr.subRouters, &router)

	return &router
}

func (rtr *Router) Methods(methods ...string) *Route {
	route := Route{
		router:  rtr,
		methods: methods,
	}

	rtr.routes = append(rtr.routes, &route)

	return &route
}

func (rt *Route) Handler(path string, handler http.Handler) {
	rt.handler = handler
	rt.path = path
}

func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !rtr.isInit {
		if err := rtr.initRoutes(); err != nil {
			return
		}
	}

	var s http.Handler = rtr.router
	if rtr.cors != nil {
		corsHandler := makeCorsHandler(rtr.cors, rtr.debugLogger)
		s = corsHandler(s)
	}

	s.ServeHTTP(w, r)
}

func (rtr *Router) initRoutes() error {
	for _, r := range rtr.routes {
		prefixedPath := r.router.prefix + r.path
		fmt.Printf("adding route for %v - %s\n", r.methods, prefixedPath)
		for _, m := range r.methods {
			if r.handler == nil {
				return fmt.Errorf("no handler for path %s", prefixedPath)
			}

			handler := r.handler
			t := len(r.router.middlewares)
			for i := t - 1; i >= 0; i-- {
				handler = r.router.middlewares[i](handler)
			}

			r.router.router.Handler(m, prefixedPath, handler)
		}
	}

	rtr.isInit = true

	for _, rs := range rtr.subRouters {
		var middlewares []middlewares.MiddlewareFunc
		middlewares = append(middlewares, rtr.middlewares...)
		middlewares = append(middlewares, rs.middlewares...)

		rs.middlewares = middlewares
		if err := rs.initRoutes(); err != nil {
			return err
		}
	}

	return nil
}

func (rtr *Router) Run(address string) error {
	if !rtr.isInit {
		if err := rtr.initRoutes(); err != nil {
			return err
		}
	}
	return http.ListenAndServe(address, rtr)
}

func GetParamsFromContext(ctx context.Context) httprouter.Params {
	return httprouter.ParamsFromContext(ctx)
}

func makeCorsHandler(config *corsConfig, debugLogger helper.DebugFieldLogger) func(http.Handler) http.Handler {
	allowMethods := config.AllowMethods
	allowHeaders := config.AllowHeaders
	exposeHeaders := config.ExposeHeaders
	maxAge := config.MaxAgeStr

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

			_ = debugLogger("event", "cors handler", "msg", fmt.Sprintf("set %s: %s", helper.HeaderAccessControlAllowOrigin, allowOrigin))

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
			_ = debugLogger("event", "cors handler", "msg", "handling pre-fligt request")
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

			_ = debugLogger("event", "cors handler", "msg", "set header for pre-flight request", "header", w.Header())
			w.WriteHeader(http.StatusOK)
		})
	}
}

// CorsAllowOrigin defines a list of origins that may access the resource.
// Optional. Default value []string{"*"}.
func CorsAllowOrigin(origins []string) CorsOption {
	return func(config *corsConfig) {
		if len(origins) > 0 {
			config.AllowOrigins = origins
		}
	}
}

// CorsAllowMethods defines a list methods allowed when accessing the resource.
// This is used in response to a preflight request.
// Optional. Default value DefaultCORSConfig.AllowMethods.
func CorsAllowMethods(methods []string) CorsOption {
	return func(config *corsConfig) {
		if len(methods) > 0 {
			config.AllowMethods = strings.Join(methods, ",")
		}
	}
}

// CorsAllowHeaders defines a list of request headers that can be used when
// making the actual request. This is in response to a preflight request.
// Optional. Default value []string{}.
func CorsAllowHeaders(headers []string) CorsOption {
	return func(config *corsConfig) {
		if len(headers) > 0 {
			config.AllowHeaders = strings.Join(headers, ",")
		}
	}
}

// CorsAllowCredentials indicates whether or not the response to the request
// can be exposed when the credentials flag is true. When used as part of
// a response to a preflight request, this indicates whether or not the
// actual request can be made using credentials.
// Optional. Default value false.
func CorsAllowCredentials(allowCredential bool) CorsOption {
	return func(config *corsConfig) {
		config.AllowCredentials = allowCredential
	}
}

// CorsMaxAge indicates how long (in seconds) the results of a preflight request
// can be cached.
// Optional. Default value 0.
func CorsMaxAge(maxAge int) CorsOption {
	return func(config *corsConfig) {
		config.MaxAge = maxAge
		config.MaxAgeStr = strconv.Itoa(maxAge)
	}
}
