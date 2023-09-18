package main

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Socket struct {
	ConnectsId int
	Connection *websocket.Conn
	User       *User
	Mutex      sync.Mutex
}

func (server *Server) Connect(connection *websocket.Conn, user *User) *Socket {
	user.Mutex.Lock()
	defer user.Mutex.Unlock()

	socket := &Socket{InsertConnect(server.SqlConn, user.Id, connection.RemoteAddr().String(), time.Now()), connection, user, sync.Mutex{}}
	return socket
}

func (server *Server) Disconnect(socket *Socket) {
	user := socket.User

	user.Mutex.Lock()
	defer user.Mutex.Unlock()

	InsertDisconnect(server.SqlConn, socket.ConnectsId, time.Now())

	socket.User = nil
	socket.ConnectsId = 0

	if user.removeSocket(socket) == 0 {
		server.removeUser(user)
	}
}
