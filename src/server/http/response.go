// response.go
package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// BoxBodyResponse is the common type for all responses.
type BoxBodyResponse struct {
	Response *http.Response
}

// ProxyResponse represents the response sent back to the client at the end of the proxying process.
type ProxyResponse struct {
	Response *http.Response
}

// NewProxyResponse creates a new ProxyResponse, wrapping the original response.
func NewProxyResponse(response *http.Response) *ProxyResponse {
	return &ProxyResponse{Response: response}
}

// IntoForwarded consumes this ProxyResponse and returns the final response that should be sent to the client.
func (pr *ProxyResponse) IntoForwarded() *http.Response {
	pr.Response.Header.Set("Server", roxyServerHeader())
	return pr.Response
}

// LocalResponse represents the HTTP response originated on this server, not obtained through a proxy process.
type LocalResponse struct{}

// Builder returns a ResponseBuilder pre-initialized with our header values.
func (lr *LocalResponse) Builder() http.Header {
	headers := http.Header{}
	headers.Set("Server", roxyServerHeader())
	return headers
}

// NotFound generates a generic HTTP 404 Not Found response.
func (lr *LocalResponse) NotFound() *http.Response {
	headers := lr.Builder()
	headers.Set("Content-Type", "text/plain")
	body := "HTTP 404 NOT FOUND"
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Header:     headers,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

// BadGateway generates a generic HTTP 502 Bad Gateway response.
func (lr *LocalResponse) BadGateway() *http.Response {
	headers := lr.Builder()
	headers.Set("Content-Type", "text/plain")
	body := "HTTP 502 BAD GATEWAY"
	return &http.Response{
		StatusCode: http.StatusBadGateway,
		Header:     headers,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

// roxyServerHeader returns the server header string.
func roxyServerHeader() string {
	return fmt.Sprintf("roxy/%s", "0.1.0") // Replace "0.1.0" with the appropriate version variable if available
}
