package session

import (
	"log/slog"
	"sync"
)

type broadcaster struct {
	clients map[string]chan int
	mu      sync.Mutex
}

func (bc *broadcaster) withLock(f func()) {
	bc.mu.Lock()
	f()
	bc.mu.Unlock()
}

var b broadcaster

func init() {
	b = broadcaster{
		clients: make(map[string]chan int),
	}
}

func ConnectClient(id string) chan int {
	var newCh chan int
	b.withLock(func() {
		if ch, exists := b.clients[id]; exists {
			// remove existing
			close(ch)
			delete(b.clients, id)
		}
		newCh = make(chan int)
		b.clients[id] = newCh
	})

	slog.Info("New Client", "session", id, "channel", newCh)
	return newCh
}

func DisconnectClient(id string) {
	b.withLock(func() {
		delete(b.clients, id)
	})
}

func Broadcast(msg int) {
	b.withLock(func() {
		for _, client := range b.clients {
			client <- msg
		}
	})
}
