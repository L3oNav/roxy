package synchronizer

import (
	"sync"

	"github.com/google/uuid"
)

type Notification int

const (
	Shutdown Notification = iota
)

// Notifier can send notifications to subscribers and await their acknowledgements.
type Notifier struct {
	mu               sync.Mutex
	subscribers      map[string]chan Notification
	acknowledgements map[string]chan struct{}
}

// Subscription allows receiving notifications from a Notifier and acknowledging them.
type Subscription struct {
	id                  string
	notificationChannel chan Notification
	acknowledgeChannel  chan struct{}
	notifier            *Notifier
}

// NewNotifier creates a new Notifier with the necessary channels.
func NewNotifier() *Notifier {
	return &Notifier{
		subscribers:      make(map[string]chan Notification),
		acknowledgements: make(map[string]chan struct{}),
	}
}

// Subscribe returns a new Subscription that can receive notifications from the Notifier.
func (n *Notifier) Subscribe() *Subscription {
	n.mu.Lock()
	defer n.mu.Unlock()

	id := uuid.New().String()
	notificationChannel := make(chan Notification, 1)
	acknowledgeChannel := make(chan struct{})

	n.subscribers[id] = notificationChannel
	n.acknowledgements[id] = acknowledgeChannel

	return &Subscription{
		id:                  id,
		notificationChannel: notificationChannel,
		acknowledgeChannel:  acknowledgeChannel,
		notifier:            n,
	}
}

// Send sends a Notification to all subscribers.
func (n *Notifier) Send(notification Notification) int {
	n.mu.Lock()
	defer n.mu.Unlock()

	count := 0
	for _, ch := range n.subscribers {
		ch <- notification
		count++
	}
	return count
}

// CollectAcknowledgements waits for all subscribers to acknowledge the last notification.
func (n *Notifier) CollectAcknowledgements() {
	n.mu.Lock()
	defer n.mu.Unlock()

	for id, ackCh := range n.acknowledgements {
		<-ackCh
		delete(n.acknowledgements, id)
		delete(n.subscribers, id)
	}
}

// ReceiveNotification reads the notifications channel for a Subscription
func (s *Subscription) ReceiveNotification() (Notification, bool) {
	notification, ok := <-s.notificationChannel
	return notification, ok
}

// AcknowledgeNotification sends an ACK on the acknowledgements channel
func (s *Subscription) AcknowledgeNotification() {
	s.acknowledgeChannel <- struct{}{}
}
