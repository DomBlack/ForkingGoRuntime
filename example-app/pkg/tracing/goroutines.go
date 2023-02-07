package tracing

import (
	"runtime"
	// unsafe allows us to use go:linkname
	_ "unsafe"

	"go.opentelemetry.io/otel/trace"
)

const spanGoRoutines = false

// goRoutineTraceData is the data that is attached to a go routine
type goRoutineTraceData struct {
	goRoutineID uint64        // The ID of the go routine
	context     *TraceContext // The trace context of the go routine
}

//go:linkname goRoutineStart runtime.tracingGStart
func goRoutineStart(pc uintptr, goRoutinueID uint64, parentTraceData *goRoutineTraceData) *goRoutineTraceData {
	if parentTraceData == nil {
		return nil
	}

	if spanGoRoutines {
		return &goRoutineTraceData{
			goRoutineID: goRoutinueID,
			context:     startSpanForOtherGoRoutine(goRoutinueID, callingFunc(pc), nil, parentTraceData.context, trace.SpanKindInternal),
		}
	} else {
		return parentTraceData
	}
}

//go:linkname goRoutineExit runtime.tracingGExit
func goRoutineExit(goRoutineID uint64, traceData *goRoutineTraceData) {
	if spanGoRoutines && traceData.context != nil {
		endSpanForOtherGoRoutine(goRoutineID, nil)
	}
}

//go:linkname goRoutineAttachData runtime.tracingAttachDataToG
func goRoutineAttachData(data *goRoutineTraceData)

//go:linkname goRoutineGetData runtime.tracingGetDataFromG
func goRoutineGetData() *goRoutineTraceData

//go:linkname goRoutineID runtime.getgoid
func goRoutineID() uint64

func goRoutineClearData() {
	goRoutineAttachData(nil)
}

func callingFunc(pc uintptr) string {
	cf := runtime.CallersFrames([]uintptr{pc})
	frame, _ := cf.Next()
	return frame.Function
}
