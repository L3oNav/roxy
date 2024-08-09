package scheduler

import (
	"fmt"
	"net"
	"roxy/src/config"
	"roxy/src/synchronizer"
)

// Scheduler interface defines the method that any load balancing algorithm
// must implement to decide which server should handle the next request.
type Scheduler interface {
	NextServer() net.Addr
}

// WeightedRoundRobin struct implements the classical Weighted Round Robin (WRR)
// algorithm for load balancing between multiple backend servers.
type WeightedRoundRobin struct {
	cycle *synchronizer.Ring[net.Addr]
}

// NewWeightedRoundRobin creates and initializes a new WeightedRoundRobin scheduler.
func NewWeightedRoundRobin(backends []config.Backend) *WeightedRoundRobin {
	cycle := []net.Addr{}

	// Interleaved WRR
	for _, backend := range backends {
		weight := backend.Weight
		addrs, err := net.ResolveTCPAddr("tcp", backend.Address)
		if err != nil {
			fmt.Errorf("Error resolving address: %v | NewWeightedRoundRobin fn", err)
		}
		for weight > 0 {
			cycle = append(cycle, addrs)
			weight--
		}
	}

	return &WeightedRoundRobin{
		cycle: synchronizer.NewRing(cycle),
	}
}

// NextServer returns the address of the server that should process the next request.
func (wrr *WeightedRoundRobin) NextServer() net.Addr {
	return wrr.cycle.NextAsOwned()
}
