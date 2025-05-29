package observ

import (
	"context"
	"net/http"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// telemetryTracer
var (
	telemetryTracer    = otel.Tracer("github.com/alifmufthi91/ecommerce-system/services/order")
	once               sync.Once
	reusablePropogator propagation.TextMapPropagator
)

// GetTracer is a helper to get the global tracker we defined.
// usually, we only create global trackers once per project or per module
func GetTracer() trace.Tracer {
	return telemetryTracer
}

// ReadTraceID is helper to read TraceID from context
func ReadTraceID(ctx context.Context) (traceID, spanID string) {
	if ctx == nil {
		return "00000000000000000000000000000000", "0000000000000000"
	}
	span := trace.SpanFromContext(ctx)
	return span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String()
}

func InjectTraceHeader(ctx context.Context, header http.Header) {
	once.Do(func() {
		reusablePropogator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	})

	reusablePropogator.Inject(ctx, propagation.HeaderCarrier(header))
}
