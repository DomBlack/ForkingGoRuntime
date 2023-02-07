package tracing

import (
	"context"
	"crypto/rand"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer trace.Tracer
	spans  = map[uint64][]spanStackEntry{}
)

type spanStackEntry struct {
	span   trace.Span
	parent *TraceContext
}

// Init initializes the tracing system under the given service name
func Init(svcName string) {
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(),
	)
	if err != nil {
		panic(err)
	}

	instanceID := [8]byte{}
	_, err = rand.Read(instanceID[:])
	if err != nil {
		panic(err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(svcName+"-svc"),
			semconv.ServiceInstanceIDKey.String(fmt.Sprintf("%x", instanceID)),
		)),
	)

	tracer = traceProvider.Tracer("example-app")
}

// startSpan starts a new span with the given name and parent span context
func startSpan(name string, remoteParent *TraceContext, kind trace.SpanKind, attrs ...attribute.KeyValue) {
	data := goRoutineGetData()
	if data == nil {
		return // not tracing this routine
	}

	traceCtx := startSpanForOtherGoRoutine(data.goRoutineID, name, remoteParent, goRoutineGetData().context, kind, attrs...)
	goRoutineGetData().context = traceCtx
}

func startSpanForOtherGoRoutine(goid uint64, name string, remoteParent *TraceContext, localParent *TraceContext, kind trace.SpanKind, attrs ...attribute.KeyValue) *TraceContext {
	if tracer == nil {
		panic("tracing not initialized")
	}

	ctx := context.Background()
	if remoteParent != nil {
		parentSpanCtx := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    remoteParent.TraceID,
			SpanID:     remoteParent.SpanID,
			TraceFlags: trace.TraceFlags(remoteParent.Flags[0]),
		})
		ctx = trace.ContextWithRemoteSpanContext(ctx, parentSpanCtx)
	} else if localParent != nil {
		parentSpanCtx := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    localParent.TraceID,
			SpanID:     localParent.SpanID,
			TraceFlags: trace.TraceFlags(localParent.Flags[0]),
		})
		ctx = trace.ContextWithSpanContext(ctx, parentSpanCtx)
	}

	_, span := tracer.Start(
		ctx, name,
		trace.WithSpanKind(kind),
		trace.WithAttributes(attrs...),
	)

	spans[goid] = append([]spanStackEntry{{
		span:   span,
		parent: localParent,
	}}, spans[goid]...)

	// Attach the span to the current goroutine
	spanCtx := span.SpanContext()
	traceCtx := &TraceContext{
		TraceID: spanCtx.TraceID(),
		SpanID:  spanCtx.SpanID(),
		Flags:   [1]byte{byte(spanCtx.TraceFlags())},
	}
	return traceCtx
}

func recordEvent(name string) {
	data := goRoutineGetData()
	if data == nil {
		// We are not tracing this goroutine
		return
	}

	entry := spans[data.goRoutineID][0]
	entry.span.AddEvent(name)
}

// endSpan ends the current span and removes it from the stack
func endSpan(err error, attrs ...attribute.KeyValue) {
	data := goRoutineGetData()
	if data == nil {
		// We are not tracing this goroutine
		return
	}

	parent := endSpanForOtherGoRoutine(data.goRoutineID, err, attrs...)

	// Restore the parent span context
	goRoutineGetData().context = parent
}

func endSpanForOtherGoRoutine(goid uint64, err error, attrs ...attribute.KeyValue) *TraceContext {
	entry := spans[goid][0]

	if len(attrs) > 0 {
		entry.span.SetAttributes(attrs...)
	}

	if err != nil {
		entry.span.RecordError(err)
	}

	entry.span.End(
		trace.WithStackTrace(true),
	)

	spans[goid] = spans[goid][1:]

	return entry.parent
}
