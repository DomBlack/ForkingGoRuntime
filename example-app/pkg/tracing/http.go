//go:build !simple

package tracing

import (
	"fmt"
	"net/http"
	"net/http/httptrace"
	// unsafe allows us to use go:linkname
	_ "unsafe"

	"go.opentelemetry.io/otel/semconv/v1.17.0/httpconv"
	"go.opentelemetry.io/otel/trace"
)

//go:linkname handlerStart net/http.tracingHandlerStart
func handlerStart(req *http.Request) {
	// Sanity check we're not already tracing, this should never happen
	// as the handlerEnd function should always be called before the next
	// request starts
	if goRoutineGetData() != nil {
		panic("go routine already has tracing data")
	}

	// Start tracing the go routine
	goRoutineAttachData(&goRoutineTraceData{goRoutineID: goRoutineID()})

	// Get the trace ID from the request header
	parentTrace, _ := ParseTraceContext(req.Header.Get(traceContextHeader))

	// Start a span
	startSpan(
		fmt.Sprintf("Handle: %s %s", req.Method, req.URL.Path),
		parentTrace,
		trace.SpanKindServer,
		httpconv.ServerRequest("", req)...,
	)
}

//go:linkname handlerEnd net/http.tracingHandlerEnd
func handlerEnd(didPanic bool) {
	// Sanity check we're tracing, this should never happen
	// as the handlerStart function should always be called before
	// the request ends
	if goRoutineGetData() == nil {
		panic("go routine has no tracing data")
	}

	var err error
	if didPanic {
		err = fmt.Errorf("panicked")
	}
	endSpan(err)
	goRoutineAttachData(nil)
}

//go:linkname startRoundTrip net/http.tracingStartRoundTrip
func startRoundTrip(req *http.Request) *http.Request {
	traceData := goRoutineGetData()
	if traceData == nil {
		// We're not tracing this request, so we don't need to do anything
		return req
	}

	startSpan(
		fmt.Sprintf("Call: %s %s", req.Method, req.URL.String()),
		nil,
		trace.SpanKindClient,
		httpconv.ClientRequest(req)...,
	)

	ctxWithTracer := httptrace.WithClientTrace(req.Context(), &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			recordEvent("Getting connection")
		},
		GotConn: func(info httptrace.GotConnInfo) {
			recordEvent("Got Connection")
		},
		PutIdleConn: nil,
		GotFirstResponseByte: func() {
			recordEvent("Received first byte")
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			recordEvent(fmt.Sprintf("DNS loookup of %s", info.Host))
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			recordEvent(fmt.Sprintf("DNS resolved to %s", info.Addrs))
		},
		ConnectStart: func(network, addr string) {
			recordEvent(fmt.Sprintf("Connecting to %s %s", network, addr))
		},
		ConnectDone: func(network, addr string, err error) {
			recordEvent(fmt.Sprintf("Connected to %s %s", network, addr))
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			recordEvent("Request sent")
		},
	})

	req.Header.Set(traceContextHeader, traceData.context.String())

	return req.WithContext(ctxWithTracer)
}

//go:linkname endRoundTrip net/http.tracingEndRoundTrip
func endRoundTrip(resp *http.Response, err error) {
	traceData := goRoutineGetData()
	if traceData == nil {
		// We're not tracing this request, so we don't need to do anything
		return
	}

	endSpan(err, httpconv.ClientResponse(resp)...)
}
