package requestid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ContextKey is
type ContextKey string

const (
	// ContextKeyReqID is the context key for RequestID
	ContextKeyReqID ContextKey = "requestID"

	// HTTPHeaderNameRequestID has the name of the header for request ID
	HTTPHeaderNameRequestID = "X-Request-ID"

	// LogFieldKeyReqID is the logfield key for RequestID
	LogFieldKeyReqID = "requestId"
)

// GetReqID will get reqID from a http request and return it as a string
func GetReqID(ctx context.Context) string {

	reqID := ctx.Value(ContextKeyReqID)

	if ret, ok := reqID.(string); ok {
		return ret
	}

	return ""
}

// AttachReqID will attach a brand new request ID to a http request
func AttachReqID(ctx context.Context) context.Context {

	reqID := uuid.New()

	return context.WithValue(ctx, ContextKeyReqID, reqID.String())
}

// Middleware will attach the reqID to the http.Request and add reqID to http header in the response
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := AttachReqID(r.Context())

		r = r.WithContext(ctx)

		h := w.Header()

		h.Add(HTTPHeaderNameRequestID, GetReqID(ctx))

		next.ServeHTTP(w, r)
	})
}

// NewLoggerFromReqIDStr creates  a *logrus.Entry that has requestID as a field. A new LogField inst will be created if log is nil
func NewLoggerFromReqIDStr(reqID string, ancestorLogger logrus.FieldLogger) logrus.FieldLogger {

	var retLogger logrus.FieldLogger = logrus.StandardLogger()

	if ancestorLogger != nil {
		retLogger = ancestorLogger
	}

	return retLogger.WithField(LogFieldKeyReqID, reqID)
}

// NewLoggerFromReqIDCtx creates a *logrus.Entry that has requestID as a field.  A new LogField inst will be created if log is ni
func NewLoggerFromReqIDCtx(ctx context.Context, ancestorLogger logrus.FieldLogger) logrus.FieldLogger {
	reqID := GetReqID(ctx)

	return NewLoggerFromReqIDStr(reqID, ancestorLogger)
}
