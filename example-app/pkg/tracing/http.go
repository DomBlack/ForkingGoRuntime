//go:build !simple

package tracing

import (
	"net/http"
	// unsafe allows us to use go:linkname
	_ "unsafe"

	"github.com/rs/zerolog/log"
)

const (
	tracingHeader = "X-Correlation-Id"
)

//go:linkname handlerStart net/http.tracingHandlerStart
func handlerStart(req *http.Request) {
	// Sanity check we're not already tracing, this should never happen
	// as the handlerEnd function should always be called before the next
	// request starts
	if goRoutineGetData() != nil {
		panic("go routine already has tracing data")
	}

	// Get the trace ID from the request header
	traceContext, err := ParseTraceContext(req.Header.Get(tracingHeader))
	if err != nil {
		// If the trace context is invalid, create a new one
		traceContext = NewTraceContext()
	} else {
		// If the trace context is valid, create s new span under it
		traceContext = traceContext.NewSpan()
	}

	traceData := &goRoutineTraceData{
		goRoutineID: goRoutineID(),
		context:     traceContext,
	}
	goRoutineAttachData(traceData)

	log.Info().Str("_trace", traceContext.String()).Str("method", req.Method).Str("path", req.URL.Path).Msg("> request started")
}

//go:linkname handlerEnd net/http.tracingHandlerEnd
func handlerEnd(didPanic bool) {
	traceData := goRoutineGetData()

	// Sanity check we're tracing, this should never happen
	// as the handlerStart function should always be called before
	// the request ends
	if traceData == nil {
		panic("go routine has no tracing data")
	}

	log.Info().Str("_trace", traceData.context.String()).Msg("> request ended")
	goRoutineClearData()
}

//go:linkname startRoundTrip net/http.tracingStartRoundTrip
func startRoundTrip(req *http.Request) {
	traceData := goRoutineGetData()
	if traceData == nil {
		// We're not tracing this request, so we don't need to do anything
		return
	}

	req.Header.Set(tracingHeader, traceData.context.String())
}

//go:linkname endRoundTrip net/http.tracingEndRoundTrip
func endRoundTrip(resp *http.Response, err error) {
	traceData := goRoutineGetData()
	if traceData == nil {
		// We're not tracing this request, so we don't need to do anything
		return
	}
}
