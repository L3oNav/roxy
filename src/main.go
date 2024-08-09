package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"roxy/src/config"
	"roxy/src/server"
	"syscall"
)

type Error string

func (e Error) Error() string { return string(e) }

func main() {
	config, err := config.NewConfig().Load("config.toml")

	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	local_server, err := server.NewMaster(config)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		fmt.Printf("Received signal: %s. Shutting down...\n", sig)
		local_server.ShutdownOn()
	}()

	go local_server.Run()
}
