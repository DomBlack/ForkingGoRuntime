package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func DoGet[Response any](ctx context.Context, port int, path ...string) (Response, error) {
	return makeClientRequest[*struct{}, Response](ctx, http.MethodGet, port, nil, path)
}

func DoPost[Response, Request any](ctx context.Context, port int, p Request, path ...string) (Response, error) {
	return makeClientRequest[Request, Response](ctx, http.MethodPost, port, p, path)
}

func DoPatch[Response, Request any](ctx context.Context, port int, p Request, path ...string) (Response, error) {
	return makeClientRequest[Request, Response](ctx, http.MethodPatch, port, p, path)
}

func DoDelete[Response any](ctx context.Context, port int, path ...string) (Response, error) {
	return makeClientRequest[*struct{}, Response](ctx, http.MethodDelete, port, nil, path)
}

func makeClientRequest[Request, Response any](ctx context.Context, method string, port int, request Request, path []string) (Response, error) {
	var zeroValue Response

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("http://localhost:%d/%s", port, joinPath(path)), nil)
	if err != nil {
		return zeroValue, err
	}

	// Set the content type
	req.Header.Set("Content-Type", "application/json")

	// Set the request body
	if method != http.MethodGet && method != http.MethodDelete {
		body, err := json.Marshal(request)
		if err != nil {
			return zeroValue, err
		}

		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	// Make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return zeroValue, err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		errToReturn := &restError{}
		if err := json.NewDecoder(res.Body).Decode(&errToReturn); err != nil {
			return zeroValue, err
		}
		return zeroValue, errToReturn
	} else {
		// Read the response
		var response response[Response]
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			return zeroValue, err
		}

		return response.Response, nil
	}
}

func joinPath(path []string) string {
	if len(path) == 0 {
		return ""
	}

	p := path[0]
	for _, s := range path[1:] {
		p += "/" + url.PathEscape(s)
	}

	return strings.TrimPrefix(p, "/")
}
