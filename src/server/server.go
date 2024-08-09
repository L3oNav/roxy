package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"roxy/src/config"
	"roxy/src/synchronizer"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

type State int

const (
	Starting State = iota
	StateListening
	StateShuttingDownPendingConnections
	StateShuttingDownDone
	StateMaxConnectionsReached
)

// ShutdownState definitions
type ShutdownState int

const (
	PendingConnections ShutdownState = iota
	Done
)

type Server struct {

	// State updates channel. Subscribers can use this to check the current
	// [`State`] of this server.
	State *atomic.Value

	// TCP listener used to accept connections.
	Listener net.Listener

	// Configuration for this server.
	Config *config.ServerConfig

	// Socket address used by this server to listen for incoming connections.
	Address string

	// [`Notifier`] object used to send notifications to tasks spawned by
	// this server.
	Notifier *synchronizer.Notifier

	// Shutdown future, this can be anything, which allows us to easily write
	// integration tests. When this future completes, the server starts the
	// shutdown process.
	Shutdown context.Context

	// Connections are limited to a maximum number. In order to allow a new
	// connection we'll have a acquire a permit from the semaphore.
	Connections *semaphore.Weighted

	mutex sync.Mutex
}

func Init(config *config.Config, replica int8) (*Server, error) {
	state := &atomic.Value{}
	state.Store(StateListening)
	var ln net.Listener
	var err error

	resolvedAddrs, err := net.ResolveTCPAddr("tcp", config.Server.LISTEN[replica])
	if err != nil {
		return nil, fmt.Errorf("failed to resolve TCP address: %w", err)
	}

	if resolvedAddrs.IP.To4() != nil {
		ln, err = net.Listen("tcp4", config.Server.LISTEN[replica])
	} else {
		ln, err = net.Listen("tcp6", config.Server.LISTEN[replica])
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	address := ln.Addr()
	notifier := synchronizer.NewNotifier()
	shutdownCtx, _ := context.WithCancel(context.Background())
	connections := semaphore.NewWeighted(int64(config.Server.MAXCONN))

	server := &Server{
		Config:      &config.Server,
		Notifier:    notifier,
		State:       state,
		Listener:    ln,
		Address:     address.String(),
		Shutdown:    shutdownCtx,
		Connections: connections,
	}

	server.State.Store(Starting)

	return server, nil

}

// Implementing the Shutdown_on method
func (s *Server) Shutdown_on() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Close the listener to stop accepting new connections
	if s.Listener != nil {
		s.Listener.Close()
	}

	// Send a shutdown notification
	s.Notifier.Send(synchronizer.Shutdown)

	// Collect acknowledgements
	s.Notifier.CollectAcknowledgements()

	log.Println("Server has been shut down gracefully.")
}

func (s *Server) Subscribe() (string, *synchronizer.Subscription) {
	return s.Address, s.Notifier.Subscribe()
}

func (s *Server) Run() error {
	config := s.Config
	state := s.State
	listener := s.Listener
	notifier := s.Notifier
	shutdown := s.Shutdown
	address := s.Address
	connections := s.Connections

	logName := address
	if config.NAME != "" {
		logName = fmt.Sprintf("%s (%s)", address, config.NAME)
	}

	config.LOGNAME = logName

	state.Store(StateListening)
	fmt.Printf("%s => Listening for requests\n", logName)

	listenerObj := &Listener{
		Config:      config,
		Connections: connections,
		Listener:    listener,
		Notifier:    notifier,
		State:       state,
	}

	errChan := make(chan error, 1)
	go func() {
		err := listenerObj.Listen()
		if err != nil {
			errChan <- fmt.Errorf("%s => Error while accepting connections: %v", logName, err)
		}
	}()

	select {
	case err := <-errChan:
		if err != nil {
			fmt.Println(err)
		}
	case <-shutdown.Done():
		fmt.Printf("%s => Received shutdown signal\n", logName)
	}

	s.Listener.Close()

	// Send shutdown notification
	if numTasks := notifier.Send(synchronizer.Shutdown); numTasks > 0 {
		if numTasks > 0 {
			fmt.Printf("%s => Can't shutdown yet, %d pending connections\n", logName, numTasks)
			state.Store(StateShuttingDownPendingConnections)
			notifier.CollectAcknowledgements()
		}
	}

	// Reset the config (simulate the unsafe drop in Rust)
	config = nil

	state.Store(StateShuttingDownDone)
	fmt.Printf("%s => Shutdown complete\n", logName)

	return nil
}

type Listener struct {
	// Server instance.
	Listener    net.Listener
	Config      *config.ServerConfig
	Notifier    *synchronizer.Notifier
	State       *atomic.Value
	Connections *semaphore.Weighted
}

func (l *Listener) Listen() error {
	for {
		if !l.Connections.TryAcquire(1) {
			fmt.Printf("%s => Reached max connections: %d\n", l.Config.LOGNAME, l.Config.MAXCONN)
			l.State.Store(StateMaxConnectionsReached)
		}

		conn, err := l.Listener.Accept()
		if err != nil {
			return err
		}

		fmt.Printf("%s => Accepted connection from %s\n", l.Config.LOGNAME, conn.RemoteAddr().String())

		go func() {
			defer l.Connections.Release(1)
			l.handleConnection(conn)
		}()
	}
}

func (l *Listener) handleConnection(conn net.Conn) {
	defer conn.Close()

	subscription := l.Notifier.Subscribe()
	defer subscription.AcknowledgeNotification()

	for {
		// Handle the connection
		notification, ok := subscription.ReceiveNotification()
		if !ok || notification == synchronizer.Shutdown {
			subscription.AcknowledgeNotification()
			break
		}
	}

	fmt.Printf("Connection from %s closed\n", conn.RemoteAddr().String())
}
