package client

import (
	"time"
)

type Manager struct {
	stopCh    chan struct{}
	clientsCh chan *Client
	clients   []*Client
}

func NewManager(timeout time.Duration, delay time.Duration, keys ...string) *Manager {
	manager := new(Manager)
	manager.stopCh = make(chan struct{})
	manager.clientsCh = make(chan *Client, len(keys))
	for i := 0; i < len(keys); i++ {
		manager.clients = append(manager.clients, NewHttpClient(timeout, delay, keys[i], manager.clientsCh))
	}
	return manager
}

func (m *Manager) Serve() error {
	for _, client := range m.clients {
		err := client.StartDelayer()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) UseClient(f func(client *Client)) {
	client := <-m.clientsCh
	f(client)
}

func (m *Manager) IsClosed() bool {
	select {
	case <-m.stopCh:
		return true
	default:
		return false
	}
}

func (m *Manager) Close() {
	if !m.IsClosed() {
		for _, client := range m.clients {
			client.Close()
		}
		close(m.stopCh)
	}
}
