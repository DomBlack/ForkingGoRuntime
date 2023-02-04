package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Get registers a GET handler
func Get[Response any](s *Server, path string, handler func(ctx *Context) (Response, error)) {
	s.router.GET(path, createHandler(s, func(ctx *Context, _ struct{}) (Response, error) {
		return handler(ctx)
	}))
}

// Post registers a POST} handler
func Post[Request, Response any](s *Server, path string, handler func(ctx *Context, r Request) (Response, error)) {
	s.router.POST(path, createHandler(s, handler))
}

// Patch registers a PATCH handler
func Patch[Request, Response any](s *Server, path string, handler func(ctx *Context, r Request) (Response, error)) {
	s.router.PATCH(path, createHandler(s, handler))
}

// Delete registers a DELETE handler
func Delete[Response any](s *Server, path string, handler func(ctx *Context) (Response, error)) {
	s.router.DELETE(path, createHandler(s, func(ctx *Context, _ struct{}) (Response, error) {
		return handler(ctx)
	}))
}

type response[Response any] struct {
	Code     int      `json:"code"`
	Response Response `json:"data"`
}

// createHandler creates a handler for the given function which
// handles marshalling and unmarshalling of the request and response
func createHandler[Request, Response any](s *Server, handler func(ctx *Context, r Request) (Response, error)) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := &Context{Context: r.Context(), headers: r.Header, params: ps}

		log := s.Log.With().Str("method", r.Method).Str("path", r.URL.Path).Logger()
		log.Debug().Msg("handling request")

		// Parse the request
		var req Request
		if r.Method != http.MethodGet && r.Method != http.MethodDelete {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Err(err).Msg("unable to parse request")
				writeError(w, BadRequest(err.Error()))
				return
			}
		}

		// Call the handler
		res, err := handler(ctx, req)
		if err != nil {
			log.Err(err).Msg("error processing request")
			writeError(w, err)
			return
		}

		wrapped := response[Response]{
			Code:     http.StatusOK,
			Response: res,
		}

		// Write the response
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(wrapped); err != nil {
			log.Err(err).Msg("unable to write response")
			writeError(w, InternalServerError(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

var (
	notFoundHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		writeError(w, NotFound(fmt.Sprintf("endpoint '%s' found", r.URL.Path)))
	}

	methodNotAllowedHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		writeError(w, MethodNotAllowed(fmt.Sprintf("method '%s' allowed for this endpoint", r.Method)))
	}

	panicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		writeError(w, InternalServerError(fmt.Sprintf("a panic occured: %v", err)))
	}
)
