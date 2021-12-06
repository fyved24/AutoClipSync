package server

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Conn struct {
	wsconn *websocket.Conn
	mu   sync.RWMutex
}

func (c *Conn) Write(message []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.wsconn.WriteMessage(websocket.TextMessage, message)
}

type ConnStore struct {
	connections map[string] *Conn
}

func (store *ConnStore) Set(key string, connection  *websocket.Conn ) bool {
	_, present := store.connections[key]
	if present {
		return false
	}
	store.connections[key] = &Conn{wsconn: connection}
	return true
}
func (store *ConnStore) Get(key string) *Conn  {
	return store.connections[key]
}
func NewConnStore() *ConnStore {
	return &ConnStore{connections: make(map[string]*Conn)}
}