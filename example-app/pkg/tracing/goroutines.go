package tracing

import (
	// unsafe allows us to use go:linkname
	_ "unsafe"
)

// goRoutineTraceData is the data that is attached to a go routine
type goRoutineTraceData struct {
	goRoutineID uint64       // The ID of the go routine
	context     TraceContext // The trace context of the go routine
}

//go:linkname goRoutineStart runtime.tracingGStart
func goRoutineStart(goRoutinueID uint64, parentTraceData *goRoutineTraceData) *goRoutineTraceData {
	return nil
}

//go:linkname goRoutineExit runtime.tracingGExit
func goRoutineExit(goRoutinueID uint64, traceData *goRoutineTraceData) {

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
