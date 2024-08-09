// proxy.go
package service

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	local_http "roxy/src/server/http"

	"golang.org/x/sync/errgroup"
)

// Forward forwards the request to the target server and returns the response sent by the target server.
func Forward(ctx context.Context, req *http.Request, targetAddr string) (*http.Response, error) {
	conn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		return new(local_http.LocalResponse).BadGateway(), nil
	}
	defer conn.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return new(local_http.LocalResponse).BadGateway(), nil
	}

	if resp.StatusCode == http.StatusSwitchingProtocols {
		clientUpgrade := req.Header.Get("Upgrade")
		serverUpgrade := resp.Header.Get("Upgrade")
		if clientUpgrade == "" || serverUpgrade == "" {
			return new(local_http.LocalResponse).BadGateway(), nil
		}
		tunnel(ctx, conn, req.Body, resp.Body)
	}

	return local_http.NewProxyResponse(resp).IntoForwarded(), nil
}

// tunnel creates a TCP tunnel for upgraded connections such as WebSockets or any other custom protocol.
func tunnel(ctx context.Context, clientConn net.Conn, clientBody io.ReadCloser, serverBody io.ReadCloser) {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer clientBody.Close()
		defer serverBody.Close()
		_, err := io.Copy(clientConn, serverBody)
		return err
	})

	eg.Go(func() error {
		defer clientBody.Close()
		defer serverBody.Close()
		_, err := io.Copy(io.Discard, serverBody)
		return err
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("Tunnel error: %v\n", err)
	}
}
