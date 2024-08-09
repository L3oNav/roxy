package service

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"roxy/src/config"
	"time"
)

type Roxy struct {
	Config     *config.Config
	ClientAddr net.Addr
	ServerAddr net.Addr
}

func NewRoxy(config *config.Config, clientAddr net.Addr, serverAddr net.Addr) *Roxy {
	return &Roxy{
		Config:     config,
		ClientAddr: clientAddr,
		ServerAddr: serverAddr,
	}
}

func (roxy *Roxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	uri := r.RequestURI
	method := r.Method

	var matchedPattern *config.Pattern
	for _, pattern := range roxy.Config.Pattern {
		if startsWith(uri, pattern.URI) {
			matchedPattern = &pattern
			break
		}
	}

	switch matchedPattern.Action.Type {
	case "forward":
		targetAddr := matchedPattern.URI
		req := r.Clone(context.Background())
		resp, err := Forward(req.Context(), req, targetAddr)
		if err != nil {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}
		copyResponse(w, resp)
	case "serve":
		// Implement file serving logic here if necessary
		http.ServeFile(w, r, *matchedPattern.Action.Serve)
	}

	logRequest(roxy.Config.Server.LOGNAME, method, uri, w, start)
}

func startsWith(str, prefix string) bool {
	return len(str) >= len(prefix) && str[:len(prefix)] == prefix
}

func copyResponse(w http.ResponseWriter, resp *http.Response) {
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func logRequest(logName, method, uri string, w http.ResponseWriter, start time.Time) {
	status := w.Header().Get("Status")
	elapsed := time.Since(start)
	fmt.Printf("%s -> %s %s %s HTTP %s %v\n", logName, method, uri, status, elapsed)
}
