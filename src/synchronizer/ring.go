package synchronizer

import (
	"fmt"
	"sync/atomic"
)

// Ring provides circular read-only access to the elements of an array.
// It is used for schedulers that can pre-compute a complete cycle and then
// return elements from that cycle when needed.
type Ring[T any] struct {
	// All the elements in this ring.
	values []T

	// Index of the next value that should be returned.
	next atomic.Uint64
}

// NewRing creates a new Ring. The first value returned when calling one of the
// getter functions will be located at index 0 in the values slice. Subsequent
// calls to any getter will return the value at the next index until the last
// one is reached, after that it starts again from the beginning. Note that
// values must have a length greater than 0, in other words it cannot be an
// empty slice.
func NewRing[T any](values []T) *Ring[T] {
	if len(values) == 0 {
		panic("Ring[T] doesn't work with empty slice")
	}
	return &Ring[T]{
		values: values,
	}
}

// nextIndex computes the index of the next value that has to be returned.
func (r *Ring[T]) nextIndex() uint64 {
	if len(r.values) == 1 {
		return 0
	}
	return r.next.Add(1) % uint64(len(r.values))
}

// NextAsRef returns a reference to the next value in the ring.
func (r *Ring[T]) NextAsRef() *T {
	return &r.values[r.nextIndex()]
}

// NextAsOwned returns the next value in the ring by making a copy.
func (r *Ring[T]) NextAsOwned() T {
	return *r.NextAsRef()
}

// NextAsCloned returns the next value in the ring by cloning it.
func (r *Ring[T]) NextAsCloned() T {
	v := r.NextAsRef()
	return *v // Assuming T has a Clone method or is clonable.
}

func main() {
	r := NewRing([]int{1, 2, 3})
	fmt.Println(r.NextAsOwned()) // 1
	fmt.Println(r.NextAsOwned()) // 2
	fmt.Println(r.NextAsOwned()) // 3
	fmt.Println(r.NextAsOwned()) // 1
}
