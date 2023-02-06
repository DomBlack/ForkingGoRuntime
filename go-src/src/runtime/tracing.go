package runtime

import (
	"unsafe"
)

// tracingGStart is called when a goroutine starts. It returns a pointer to a
// the parent routines trace data, and is expected to return a pointer to the
// new routines trace data.
//
// If the parent Go routine has no tracing data, this won't get called
// thus it is always safe to assume parent is non-nil
func tracingGStart(goRoutinueID uint64, parentTraceData unsafe.Pointer) unsafe.Pointer

// tracingGExit is called when a goroutine exits. It is passed a pointer to the
// trace data of the exiting goroutine.
//
// If no tracing data is available, this won't get called
// thus it is always safe to assume parent is non-nil
func tracingGExit(goRoutinueID uint64, traceData unsafe.Pointer)

// tracingAttachDataToG attaches the given data to the current goroutine.
func tracingAttachDataToG(data unsafe.Pointer) {
	getg().traceData = data
}

// tracingGetDataFromG returns the tracing data attached to the current goroutine.
func tracingGetDataFromG() unsafe.Pointer {
	return getg().traceData
}

// getgoid returns the ID of the current goroutine.
func getgoid() uint64 {
	return getg().goid
}
