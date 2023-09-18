package main

import (
	"database/sql"
	"math/rand"
	"sync"
)

type Random struct {
	RandSeed      rand.Source
	RandSeedMutex *sync.Mutex
}

type Server struct {
	SqlConn *sql.DB
	Rand    *Random

	users      map[int]*User
	usersMutex *sync.Mutex
}

func Create() *Server {
	return &Server{
		SqlConn:    ConnectToLocalDb("db username", "db password", "db name"),
		Rand:       &Random{rand.NewSource(rand.Int63()), &sync.Mutex{}},
		users:      make(map[int]*User),
		usersMutex: &sync.Mutex{},
	}
}
