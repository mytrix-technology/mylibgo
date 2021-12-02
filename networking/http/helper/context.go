package helper

import (
	"context"
	"net/http"
	"strconv"
)

// PopulateRequestContext is a RequestFunc that populates several values into
// the context from the HTTP request. Those values may be extracted using the
// corresponding ContextKey type in this package.
func PopulateRequestContext(ctx context.Context, r *http.Request) context.Context {
	for k, v := range map[contextKey]string{
		ContextKeyRequestMethod:          r.Method,
		ContextKeyRequestURI:             r.RequestURI,
		ContextKeyRequestPath:            r.URL.Path,
		ContextKeyRequestProto:           r.Proto,
		ContextKeyRequestHost:            r.Host,
		ContextKeyRequestRemoteAddr:      r.RemoteAddr,
		ContextKeyRequestXForwardedFor:   r.Header.Get("X-Forwarded-For"),
		ContextKeyRequestXForwardedProto: r.Header.Get("X-Forwarded-Proto"),
		ContextKeyRequestAuthorization:   r.Header.Get("Authorization"),
		ContextKeyRequestReferer:         r.Header.Get("Referer"),
		ContextKeyRequestUserAgent:       r.Header.Get("User-Agent"),
		ContextKeyRequestXRequestID:      r.Header.Get("X-Request-Id"),
		ContextKeyRequestAccept:          r.Header.Get("Accept"),
		ContextKeyRequestXTraceID: r.Header.Get("X-Trace-Id"),
		ContextKeyRequestDatetime: r.Header.Get("datetime"),
		ContextKeyRequestSignature: r.Header.Get("signature"),
		ContextKeyUTUserID: r.Header.Get("UT-UserID"),
		ContextKeyUTToken: r.Header.Get("UT_Token"),
	} {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

type contextKey int

const (
	// ContextKeyRequestMethod is populated in the context by
	// PopulateRequestContext. Its value is r.Method.
	ContextKeyRequestMethod contextKey = iota

	// ContextKeyRequestURI is populated in the context by
	// PopulateRequestContext. Its value is r.RequestURI.
	ContextKeyRequestURI

	// ContextKeyRequestPath is populated in the context by
	// PopulateRequestContext. Its value is r.URL.Path.
	ContextKeyRequestPath

	// ContextKeyRequestProto is populated in the context by
	// PopulateRequestContext. Its value is r.Proto.
	ContextKeyRequestProto

	// ContextKeyRequestHost is populated in the context by
	// PopulateRequestContext. Its value is r.Host.
	ContextKeyRequestHost

	// ContextKeyRequestRemoteAddr is populated in the context by
	// PopulateRequestContext. Its value is r.RemoteAddr.
	ContextKeyRequestRemoteAddr

	// ContextKeyRequestXForwardedFor is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("X-Forwarded-For").
	ContextKeyRequestXForwardedFor

	// ContextKeyRequestXForwardedProto is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("X-Forwarded-Proto").
	ContextKeyRequestXForwardedProto

	// ContextKeyRequestAuthorization is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("Authorization").
	ContextKeyRequestAuthorization

	// ContextKeyRequestReferer is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("Referer").
	ContextKeyRequestReferer

	// ContextKeyRequestUserAgent is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("User-Agent").
	ContextKeyRequestUserAgent

	// ContextKeyRequestXRequestID is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("X-Request-Id").
	ContextKeyRequestXRequestID

	// ContextKeyRequestAccept is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("Accept").
	ContextKeyRequestAccept

	// ContextKeyResponseHeaders is populated in the context whenever a
	// ServerFinalizerFunc is specified. Its value is of type http.Header, and
	// is captured only once the entire response has been written.
	ContextKeyResponseHeaders

	// ContextKeyResponseSize is populated in the context whenever a
	// ServerFinalizerFunc is specified. Its value is of type int64.
	ContextKeyResponseSize

	// ContextKeyRequestXTraceID is populated in the context by
	// PopulateRequestContext. Its value is r.Header.Get("X-Trace-Id").
	ContextKeyRequestXTraceID

	ContextKeyRequestDatetime
	ContextKeyRequestSignature

	ContextKeyUTUserID
	ContextKeyUTToken
)

func ReqIDFromContext(ctx context.Context) (string, bool) {
	reqID, ok := ctx.Value(ContextKeyRequestXRequestID).(string)
	if !ok {
		return "", false
	}
	return reqID, ok
}

func TraceIDFromContext(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(ContextKeyRequestXTraceID).(string)
	if !ok {
		return "", false
	}
	return traceID, ok
}

func UTUserIDFromContext(ctx context.Context) (int64, bool) {
	userIDInt := int64(0)
	userID, ok := ctx.Value(ContextKeyUTUserID).(string)
	if ok {
		if userID != "" {
			in, err := strconv.ParseInt(userID, 10, 64)
			if err == nil {
				userIDInt = in
			}
		}
	}

	return userIDInt, true
}

func UTTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(ContextKeyUTToken).(string)
	if !ok {
		return "", false
	}
	return token, ok
}