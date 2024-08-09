package server

import (
	"context"
	"fmt"
	"roxy/src/config"
	"roxy/src/synchronizer"
	"sync"
)

type Master struct {
	Servers        []*Server
	States         []StateInfo
	Shutdown       context.Context
	ShutdownCancel context.CancelFunc
}

type StateInfo struct {
	Address  string
	StateSub *synchronizer.Subscription
}

func NewMaster(config *config.Config) (*Master, error) {
	var servers []*Server
	var states []StateInfo
	ctx, cancel := context.WithCancel(context.Background())

	for index, _ := range config.Server.LISTEN {
		server, err := Init(config, int8(index))
		if err != nil {
			return nil, err
		}
		address, sub := server.Subscribe()
		states = append(states, StateInfo{Address: address, StateSub: sub})
		servers = append(servers, server)
	}

	return &Master{
		Servers:        servers,
		States:         states,
		Shutdown:       ctx,
		ShutdownCancel: cancel,
	}, nil
}

func (m *Master) ShutdownOn() {
	m.ShutdownCancel()
}

func (m *Master) Run() error {
	var wg sync.WaitGroup

	for _, server := range m.Servers {
		wg.Add(1)
		go func(s *Server) {
			defer wg.Done()
			if err := s.Run(); err != nil {
				fmt.Printf("Error running server: %v\n", err)
			}
		}(server)
	}

	<-m.Shutdown.Done()
	fmt.Println("Master => Sending shutdown signal to all servers")
	for _, server := range m.Servers {
		server.Shutdown_on()
	}

	wg.Wait()
	fmt.Println("Master => All servers have shut down")
	return nil
}

func (m *Master) Sockets() []string {
	var addresses []string
	for _, state := range m.States {
		addresses = append(addresses, state.Address)
	}
	return addresses
}
