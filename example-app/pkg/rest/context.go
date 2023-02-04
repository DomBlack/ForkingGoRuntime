package rest

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Context is a wrapper around the standard context.Context
// which includes the request headers and path parameters
// on server requests
type Context struct {
	context.Context

	headers http.Header
	params  httprouter.Params
}

// Param returns the path parameter with the given name from the request
func (c *Context) Param(name string) string {
	return c.params.ByName(name)
}

// Header returns the header with the given name from the request
func (c *Context) Header(name string) string {
	return c.headers.Get(name)
}
