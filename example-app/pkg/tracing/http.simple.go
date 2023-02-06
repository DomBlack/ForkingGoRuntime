//go:build simple

package tracing

import (
	"net/http"
	// unsafe allows us to use go:linkname
	_ "unsafe"

	"github.com/rs/zerolog/log"
)

//go:linkname handlerStart net/http.tracingHandlerStart
func handlerStart(req *http.Request) {
	// Sanity check we're not already tracing, this should never happen
	// as the handlerEnd function should always be called before the next
	// request starts
	if goRoutineGetData() != nil {
		panic("go routine already has tracing data")
	}

	traceData := &goRoutineTraceData{
		goRoutineID: goRoutineID(),
	}
	goRoutineAttachData(traceData)

	log.Info().Str("method", req.Method).Str("url", req.URL.String()).Uint64("goid", traceData.goRoutineID).Msg("> request started")
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

	log.Info().Uint64("goid", traceData.goRoutineID).Msg("> request finished")
	goRoutineClearData()
}

//go:linkname startRoundTrip net/http.tracingStartRoundTrip
func startRoundTrip(req *http.Request) {
	if traceData := goRoutineGetData(); traceData != nil {
		log.Info().Uint64("goid", traceData.goRoutineID).Str("method", req.Method).Str("url", req.String()).Msg("> starting round trip")
	}
}

//go:linkname endRoundTrip net/http.tracingEndRoundTrip
func endRoundTrip(resp *http.Response, err error) {
	if traceData := goRoutineGetData(); traceData != nil {
		log.Info().Uint64("goid", traceData.goRoutineID).Err(err).Str("status", resp.Status).Msg("> ended round trip")
	}
}
