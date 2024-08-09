package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	port := flag.String("port", "3312", "Port on which the server will run")
	flag.Parse()

	addr := fmt.Sprintf("localhost:%s", *port)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go func(c net.Conn) {
			defer c.Close()
			c.Write([]byte("Hello, World!"))
		}(conn)
	}
}
