package main

import "sync"

type User struct {
	Id   int
	Name string

	Sockets       map[int]*Socket
	SocketCounter int

	Mutex sync.Mutex
}

func (user *User) addSocket(socket *Socket) int {

	user.Sockets[socket.ConnectsId] = socket
	user.SocketCounter++
	return user.SocketCounter
}

func (user *User) removeSocket(socket *Socket) int {

	delete(user.Sockets, socket.ConnectsId)
	user.SocketCounter--
	return user.SocketCounter
}

func (server *Server) removeUser(user *User) {
	server.usersMutex.Lock()
	defer server.usersMutex.Unlock()

	delete(server.users, user.Id)
}

func (server *Server) getOrAddUser(userId int) *User {

	server.usersMutex.Lock()
	defer server.usersMutex.Unlock()

	if user := server.users[userId]; user != nil {
		return user
	} else {
		user := &User{
			Id:   userId,
			Name: SelectUsernameByUserId(server.SqlConn, userId),

			Sockets: map[int]*Socket{},
			Mutex:   sync.Mutex{},
		}
		server.users[userId] = user
		return user
	}
}
