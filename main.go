package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/google/uuid"
)

func main() {
	// Create a map of connections
	space := make(map[string]*Connection)
	rooms := []Room{
		Room{name: "coolblue"},
		Room{name: "bottleup"},
	}
	// Lets make some rooms

	fmt.Println("Listening on TCP 5000")

	l, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	// Close the connection
	defer l.Close()

	for {

		// accept connections to our tcp server
		// lets see if we can also create a struct for it?
		// key can be some unique uuid
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		u := uuid.NewString()
		// on initial connection, no room
		c := Connection{connection: conn, id: u, room: "None"}
		// make a uuid

		// add to map
		space[u] = &c
		fmt.Println(space)

		go logic(*space[u], &rooms)

	}

}

// This will probably be the main connection handler with the users I think
func logic(c Connection, r *[]Room) {
	// finally close the connection once the handle connection is done
	defer c.connection.Close()

	//first things first, lets possibly let the user know its uuid
	sendBackToClient(c.connection, "your uuid: "+c.id+"\n")
	fmt.Println(*r)

	// So here, i think we should ask the user to join a room...lets have 2 rooms for now
	for i := 0; i < len(*r); i++ {

		number := (strconv.Itoa(i + 1))

		// Interestingly...dereference a pointer first before indexing :p i.e (*p)[0]
		sendBackToClient(c.connection, number+": "+(*r)[i].name+"\n")

	}
	sendBackToClient(c.connection, "any room preference? type : room={{room number}} to proceed\n>")

	// r.joinRoom(&c)

	// fmt.Printf("Room: %v", *r)

	// The following is basically reading input from user...good channel / goroutine usecase?
	// i want people in the same room to chat
	// for {
	// 	input, err := bufio.NewReader(c.connection).ReadString('\n')
	// 	if err != nil {
	// 		fmt.Printf("Error: %v", err)
	// 		break
	// 	}
	// 	sendBackToClient(c.connection, ">you are not in a room so nobody read this...\n>")
	// 	fmt.Println(input)
	// }
	readFromClient(&c, r)

}

func readFromClient(c *Connection, r *[]Room) {
	// read input from client and do something i dont know what yet
	for {
		input, err := bufio.NewReader(c.connection).ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v", err)
			break
		}

		// check if user wants to join a room
		if input == "room=1\n" {
			sendBackToClient(c.connection, ">joining room 1\n>")
			c.room = "coolblue"
			(*r)[0].joinRoom(c)

		}
		// broadcast user message if in a specific room? lets func that
		if c.room == "coolblue" {
			broadcast(c, &(*r)[0], input)

		}
		if c.room == "bottleup" {
			broadcast(c, &(*r)[0], input)

		}

		if c.room == "none" {
			sendBackToClient(c.connection, ">you are not in a room so nobody read this...\n>")
		}

		// fmt.Println(input)
	}

}

func broadcast(c *Connection, r *Room, s string) {

	for i := 0; i < len((*r).users); i++ {

		// lets send some data back over the connection if NOT current user
		if (*r).users[i].id != c.id {
			if _, err := io.WriteString(c.connection, s); err != nil {
				fmt.Printf("Error: %v", err)
			}
		}
	}

}

func sendBackToClient(c net.Conn, s string) {

	// lets send some data back over the connection
	if _, err := io.WriteString(c, s); err != nil {
		fmt.Printf("Error: %v", err)
	}

}
