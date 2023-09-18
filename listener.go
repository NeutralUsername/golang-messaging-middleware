package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Start(server *Server) {
	http.HandleFunc("/", frontEndHandler("*frontend path*"))
	http.HandleFunc("/ws", wsHandler(server))
	log.Println("server started on 443/80")
	go func() {
		err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://playnodge.com", http.StatusMovedPermanently)
		}))
		if err != nil {
			log.Fatal(err)
		}
	}()
	err := http.ListenAndServeTLS(":443", "*cert path*", "*key path*", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func frontEndHandler(frontEndPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir(frontEndPath)).ServeHTTP(w, r)
	}
}

func wsHandler(server *Server) func(http.ResponseWriter, *http.Request) {
	return func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		if connection, err := upgrader.Upgrade(httpResponseWriter, httpRequest, nil); err == nil {
			ListenAndServe(server, server.Connect(connection, server.AuthenticateRequest(httpRequest)))
		}
	}
}

func (server *Server) AuthenticateRequest(httpRequest *http.Request) *User {
	id := 0
	if nameCookie, err1 := httpRequest.Cookie("username"); err1 == nil {
		if passwordCookie, err2 := httpRequest.Cookie("password"); err2 == nil {
			if userId := SelectUserIdByUsernameAndPassword(server.SqlConn, nameCookie.Value, passwordCookie.Value); userId != 0 {
				id = userId
			}
		}
	}
	if id == 0 {
		for id = InsertUser(server.SqlConn, server.Rand.GenerateRandomString(10, VALID_NAME_CHARS),
			Sha256string(server.Rand.GenerateRandomString(15, VALID_NAME_CHARS)), USER_POWER_GUEST, time.Now()); id == 0; {
		}
	}
	return server.getOrAddUser(id)
}

func Sha256string(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func (rand *Random) GenerateRandomString(length int, validChars string) string {
	str := ""
	for i := 0; i < length; i++ {
		rand.RandSeedMutex.Lock()
		str += string(validChars[rand.RandSeed.Int63()%int64(len(validChars))])
		rand.RandSeedMutex.Unlock()
	}
	return str
}
