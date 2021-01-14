package tracing

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// LogError logs an error to an opentracing span
func LogError(span opentracing.Span, message string) {
	span.SetTag("error", true)
	span.LogFields(
		log.String("event", "error"),
		log.String("message", message),
	)
}
