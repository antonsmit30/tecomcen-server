package main

import (
	"fmt"
	"net"
)

type Connection struct {
	connection net.Conn
	id         string
	room       string
	name       string
}

type Room struct {
	name     string
	users    []Connection
	messages chan string
}

// methods

// joins a room
func (r *Room) joinRoom(c *Connection) {

	fmt.Println("Joining room")
	r.users = append(r.users, *c)

}
