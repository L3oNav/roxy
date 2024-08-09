// body.go
package http

import (
	"bytes"
	"io"
	"net/http"
)

// Full returns a single chunk body.
func Full(chunk string) *http.Request {
	body := io.NopCloser(bytes.NewReader([]byte(chunk)))
	req, _ := http.NewRequest("GET", "/", body)
	return req
}

// Empty returns an empty body.
func Empty() *http.Request {
	req, _ := http.NewRequest("GET", "/", nil)
	return req
}
