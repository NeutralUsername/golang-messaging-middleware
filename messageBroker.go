package main

import (
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func ListenAndServe(server *Server, socket *Socket) {
	for {
		_, message, err := socket.Connection.ReadMessage()
		if err != nil {
			break
		}
		messageHandler(server, socket, strings.Split(string(message), WS_MSG_DELIMITER))
	}
	socket.Connection.Close()
	server.Disconnect(socket)
}

func messageHandler(server *Server, socket *Socket, segments []string) {
	socket.User.Mutex.Lock()
	defer socket.User.Mutex.Unlock()

	InsertActivity(server.SqlConn, socket.ConnectsId, segments[0], strings.Join(segments[1:], WS_MSG_DELIMITER), time.Now())

	switch segments[0] {
	case "ping":
		MessageSocket(socket, ConstructMessage("pong", ""))
	}
}

func ConstructMessage(msgType string, msgData ...string) []byte {
	msg := msgType
	for _, data := range msgData {
		msg += WS_MSG_DELIMITER + data
	}
	return []byte(msg)
}

func MessageSocket(socket *Socket, message_data []byte) {
	socket.Mutex.Lock()
	defer socket.Mutex.Unlock()
	if err := socket.Connection.WriteMessage(websocket.TextMessage, message_data); err != nil {
		println("Failed to send message to socket with ip " + socket.Connection.RemoteAddr().String() + ": " + err.Error())
	}
}
