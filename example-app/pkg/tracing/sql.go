package tracing

import (
	"fmt"
	_ "unsafe"

	"go.opentelemetry.io/otel/trace"
)

//go:linkname sqlQueryStart database/sql.tracingQueryStart
func sqlQueryStart(query string) {
	traceData := goRoutineGetData()
	if traceData == nil {
		// We're not tracing this request, so we don't need to do anything
		return
	}

	startSpan(
		fmt.Sprintf("%s", query),
		nil,
		trace.SpanKindInternal,
	)
}

//go:linkname sqlQueryEnd database/sql.tracingQueryEnd
func sqlQueryEnd(err error) {
	traceData := goRoutineGetData()
	if traceData == nil {
		// We're not tracing this request, so we don't need to do anything
		return
	}

	endSpan(err)
}
