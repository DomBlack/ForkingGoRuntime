package tracing

import (
	"fmt"
	"net/http"
	_ "unsafe"
)

//go:linkname handlerStart net/http.tracingHandlerStart
func handlerStart(req *http.Request) {
	traceData := &goRoutineTraceData{traceID: goRoutineID()}
	goRoutineAttachData(traceData)

	fmt.Printf("[%d] Request Started: %s %s\n", traceData.traceID, req.Method, req.URL.String())
}

//go:linkname handlerEnd net/http.tracingHandlerEnd
func handlerEnd(didPanic bool) {
	fmt.Printf("[%d] Requested Finished\n", goRoutineGetData().traceID)
	goRoutineAttachData(nil)
}

//go:linkname startRoundTrip net/http.tracingStartRoundTrip
func startRoundTrip(req *http.Request) {
	if traceData := goRoutineGetData(); traceData != nil {
		fmt.Printf("[%d] Calling %s %s\n", traceData.traceID, req.Method, req.URL.String())
	}
}

//go:linkname endRoundTrip net/http.tracingEndRoundTrip
func endRoundTrip(resp *http.Response, err error) {
	if traceData := goRoutineGetData(); traceData != nil {
		fmt.Printf("[%d] Finished Call with status: %s\n", traceData.traceID, resp.Status)
	}
}
