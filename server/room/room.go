package room

import (
	"net"
	"sync"
)

type Room struct {
	name    string
	clients map[net.Conn]bool
	mtx     sync.Mutex
}

var (
	rooms   = make(map[string]*Room)
	roomsMtx sync.Mutex
)

// GET sebuah room yang sudah ada atau membuat room yang baru
func GetOrCreateRoom(name string) *Room {
	roomsMtx.Lock()
	defer roomsMtx.Unlock()

	if r, exists := rooms[name]; exists {
		return r
	}

	newRoom := &Room{
		name:    name,
		clients: make(map[net.Conn]bool),
	}
	rooms[name] = newRoom
	return newRoom
}

// ADD client ke sebuah room
func (r *Room) AddClient(conn net.Conn) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.clients[conn] = true
}

// REMOVE client dari sebuah room
func (r *Room) RemoveClient(conn net.Conn) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	delete(r.clients, conn)

	// Kalau room kosong, delete room tersebut
	if len(r.clients) == 0 {
		roomsMtx.Lock()
		delete(rooms, r.name)
		roomsMtx.Unlock()
	}
}

// Broadcast message ke semua room kecuali server
func (r *Room) Broadcast(sender net.Conn, message string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for conn := range r.clients {
		if conn != sender {
			conn.Write([]byte(message))
		}
	}
}
