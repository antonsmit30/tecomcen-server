package main

import "net"

type Connection struct {
	connection net.Conn
	id         string
	room       string
}

type Room struct {
	name  string
	users []Connection
}

// methods

// joins a room
func (r *Room) joinRoom(c *Connection) {

	r.users = append(r.users, *c)

}
