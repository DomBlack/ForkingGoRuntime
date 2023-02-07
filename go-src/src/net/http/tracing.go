package http

import (
	_ "unsafe"
)

// tracingHandlerStart is called when a HTTP request starts.
func tracingHandlerStart(req *Request)

// tracingHandlerEnd is called when a HTTP request ends.
//
// If the handler panicked, didPanic will be true
// otherwise it will be false.
func tracingHandlerEnd(didPanic bool)

// tracingStartRoundTrip is called when a HTTP request starts.
func tracingStartRoundTrip(req *Request) *Request

// tracingEndRoundTrip is called when a HTTP request ends.
func tracingEndRoundTrip(resp *Response, err error)
