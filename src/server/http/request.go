package http

import (
	"fmt"
	"net"
	"net/http"
)

// ProxyRequest represents a request received by the proxy from a client.
type ProxyRequest struct {
	Request    *http.Request
	ClientAddr net.Addr
	ServerAddr net.Addr
	ProxyID    *string
}

// NewProxyRequest creates a new ProxyRequest.
func NewProxyRequest(request *http.Request, clientAddr, serverAddr net.Addr, proxyID *string) *ProxyRequest {
	return &ProxyRequest{
		Request:    request,
		ClientAddr: clientAddr,
		ServerAddr: serverAddr,
		ProxyID:    proxyID,
	}
}

// Headers returns the headers of the request.
func (pr *ProxyRequest) Headers() http.Header {
	return pr.Request.Header
}

// IntoForwarded consumes the ProxyRequest and returns an http.Request that contains a valid HTTP forwarded header.
func (pr *ProxyRequest) IntoForwarded() *http.Request {
	host := pr.Request.Host
	if host == "" {
		host = pr.ServerAddr.String()
	}

	by := pr.ServerAddr.String()
	if pr.ProxyID != nil {
		by = *pr.ProxyID
	}

	forwarded := fmt.Sprintf("for=%s;by=%s;host=%s", pr.ClientAddr.String(), by, host)

	if existingForwarded := pr.Request.Header.Get("Forwarded"); existingForwarded != "" {
		forwarded = fmt.Sprintf("%s, %s", existingForwarded, forwarded)
	}

	pr.Request.Header.Set("Forwarded", forwarded)

	return pr.Request
}
