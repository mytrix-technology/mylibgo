package middlewares

func NewHttpRequestIDInjectorMiddleware() MiddlewareFunc {
	return HttpRequestIDInjectorMiddleware
}