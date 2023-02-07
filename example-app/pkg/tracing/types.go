package tracing

import (
	"errors"
	"fmt"
)

const traceContextHeader = "traceparent"

type TraceContext struct {
	TraceID [16]byte
	SpanID  [8]byte
	Flags   [1]byte
}

func ParseTraceContext(s string) (*TraceContext, error) {
	if s == "" {
		return nil, errors.New("trace context is empty")
	}

	traceID := make([]byte, 16)
	spanID := make([]byte, 8)
	flags := make([]byte, 1)

	_, err := fmt.Sscanf(s, "00-%x-%x-%x", &traceID, &spanID, &flags)
	if err != nil {
		return nil, err
	}

	rtn := &TraceContext{}
	copy(rtn.TraceID[:], traceID)
	copy(rtn.SpanID[:], spanID)
	copy(rtn.Flags[:], flags)
	return rtn, nil
}

func (tc TraceContext) String() string {
	return fmt.Sprintf("00-%x-%x-%x", tc.TraceID, tc.SpanID, tc.Flags)
}
