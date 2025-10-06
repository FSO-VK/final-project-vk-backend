// Package http provides utility functions for handling HTTP responses.
package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
	"github.com/valyala/fasthttp"
)

// ErrResponseIsNil is returned when a response object is unexpectedly nil.
var ErrResponseIsNil = errors.New("response is nil")

// WriteJSON writes encoded JSON response to the io.Writer.
func EncodeJSON[T any](w io.Writer, response *api.Response[T]) error {
	if response == nil {
		return ErrResponseIsNil
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}

// NetHTTPWriteJSON writes encoded JSON response to the http.ResponseWriter.
func NetHTTPWriteJSON[T any](w http.ResponseWriter, response *api.Response[T]) error {
	w.Header().Set("Content-Type", "application/json")

	return EncodeJSON(w, response)
}

// FastHTTPWriteJSON writes encoded JSON response to the fasthttp.RequestCtx.
func FastHTTPWriteJSON[T any](ctx *fasthttp.RequestCtx, response *api.Response[T]) error {
	ctx.SetContentType("application/json")

	return EncodeJSON(ctx, response)
}
