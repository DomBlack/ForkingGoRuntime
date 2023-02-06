package tracing

import (
	"crypto/rand"
	"errors"
	"fmt"
)

const traceContextHeader = "traceparent"

type TraceContext struct {
	TraceID    [16]byte
	ParentSpan [8]byte
	SpanID     [8]byte
	Flags      [1]byte
}

func NewTraceContext() (rtn TraceContext) {
	_, err := rand.Read(rtn.TraceID[:])
	if err != nil {
		panic("unable to generate trace ID")
	}
	_, err = rand.Read(rtn.SpanID[:])
	if err != nil {
		panic("unable to generate span ID")
	}
	return rtn
}

func ParseTraceContext(s string) (rtn TraceContext, err error) {
	if s == "" {
		return rtn, errors.New("trace context is empty")
	}

	traceID := make([]byte, 16)
	spanID := make([]byte, 8)
	flags := make([]byte, 1)

	_, err = fmt.Sscanf(s, "00-%x-%x-%x", &traceID, &spanID, &flags)
	if err != nil {
		return rtn, err
	}

	copy(rtn.TraceID[:], traceID)
	copy(rtn.SpanID[:], spanID)
	copy(rtn.Flags[:], flags)
	return rtn, nil
}

func (tc TraceContext) NewSpan() TraceContext {
	rtn := tc
	rtn.ParentSpan = tc.SpanID
	_, err := rand.Read(rtn.SpanID[:])
	if err != nil {
		panic("unable to generate span ID")
	}
	return rtn
}

func (tc TraceContext) String() string {
	return fmt.Sprintf("00-%x-%x-%x", tc.TraceID, tc.SpanID, tc.Flags)
}
