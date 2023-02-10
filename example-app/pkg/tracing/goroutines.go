package tracing

import _ "unsafe" // unsafe allows us to use go:linkname

// goRoutineTraceData is the data that is attached to a go routine
type goRoutineTraceData struct {
	traceID uint64
}

//go:linkname goRoutineStart runtime.tracingGStart
func goRoutineStart(goRoutineID uint64, parentTraceData *goRoutineTraceData) *goRoutineTraceData {
	return parentTraceData
}

//go:linkname goRoutineExit runtime.tracingGExit
func goRoutineExit(goRoutineID uint64, traceData *goRoutineTraceData) {}

//go:linkname goRoutineAttachData runtime.tracingAttachDataToG
func goRoutineAttachData(data *goRoutineTraceData)

//go:linkname goRoutineGetData runtime.tracingGetDataFromG
func goRoutineGetData() *goRoutineTraceData

//go:linkname goRoutineID runtime.tracingGetGOID
func goRoutineID() uint64
