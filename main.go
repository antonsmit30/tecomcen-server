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
	// Lets make some rooms
	space := make(map[string]*Connection)
	rooms := []Room{
		Room{name: "coolblue", messages: make(chan string)},
		Room{name: "bottleup", messages: make(chan string)},
	}
	cb_ch := make(chan (string), 10)
	bu_ch := make(chan (string), 10)
	// done := make(chan (bool), 10)

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

		go logic(*space[u], &rooms, cb_ch, bu_ch)

		// select {
		// case <-cb_ch:
		// 	fmt.Println("Top level cb channel msg received. send to entire room")
		// 	for i := 0; i < len(rooms[0].users); i++ {

		// 		go writeToClient(&rooms[0].users[i], <-cb_ch, done)

		// 	}

		// case <-bu_ch:
		// 	fmt.Println("top level bu channel msg received. we should send to entire room")
		// 	for i := 0; i < len(rooms[1].users); i++ {

		// 		go writeToClient(&rooms[1].users[i], <-bu_ch, done)

		// 	}
		// default:
		// 	fmt.Println("Default")
		// }

		// Maybe...our rooms, should message the users over here?

	}

}

// This will probably be the main connection handler with the user I think
func logic(c Connection, r *[]Room, cb chan (string), bu chan (string)) {

	local_channel := make(chan (string))

	// finally close the connection once the handle connection is done
	// defer c.connection.Close()

	// ask for room to join
	//first things first, lets possibly let the user know its uuid
	sendBackToClient(c.connection, "your uuid: "+c.id+"\n")
	fmt.Println(*r)
	sendBackToClient(c.connection, "type your username for chat: ")
	username, err := bufio.NewReader(c.connection).ReadString('\n')
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	c.name = username

	// So here, i think we should ask the user to join a room...lets have 2 rooms for now
	for i := 0; i < len(*r); i++ {

		number := (strconv.Itoa(i + 1))

		// Interestingly...dereference a pointer first before indexing :p i.e (*p)[0]
		sendBackToClient(c.connection, number+": "+(*r)[i].name+"\n")

	}
	sendBackToClient(c.connection, "any room preference? type : room={{room number}} to proceed\n>")
	input, err := bufio.NewReader(c.connection).ReadString('\n')
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	// check if user wants to join a room
	if input == "1\n" {
		sendBackToClient(c.connection, ">joining room 1\n>")
		c.room = "coolblue"
		(*r)[0].joinRoom(&c)

	} else {
		sendBackToClient(c.connection, ">You are not broadcasting to any specific room. Exiting\n>")
		return
	}
	// infinite loop to read and write messages
	for {

		for i := 0; i < len(*r); i++ {

			if (*r)[i].name == c.room {
				// fmt.Println("go readingfromclient")

				go readFromClient(&c, &(*r)[i], local_channel)
				// select {

				// msg := <-(*r)[i].messages
				// broadcast the message to each user :p
				// TODO: optimize this shit
				// for j := 0; j < len((*r)[i].users); j++ {

				// 	if (*r)[i].users[j].id != c.id {
				// 		fmt.Printf("writing to user: %v", msg)
				// 		go writeToClient(&c, msg)
				// 	}

				// }

				// }

			}

		}

		select {
		// case cb_msg := <-(*r)[0].messages:
		case cb_msg := <-local_channel:
			fmt.Printf("cb message received")
			fmt.Printf("the message received: %v", cb_msg)
			// cb <- cb_msg
			for i := 0; i < len((*r)[0].users); i++ {

				// here...we only want to broadcast the message to other users except
				if (*r)[0].users[i].id != c.id {
					go writeToClient(&(*r)[0].users[i], cb_msg, &c.name)
				}

			}

		case bu_msg := <-(*r)[1].messages:
			fmt.Printf("bu message received")
			fmt.Printf("the message received: %v", bu_msg)
			// bu <- bu_msg
			for i := 0; i < len((*r)[1].users); i++ {

				if (*r)[1].users[i].id != c.id {
					go writeToClient(&(*r)[1].users[i], bu_msg, &c.name)
				}

			}

		}

		// grab data back from client and send to room if applicable
		// we basically need to loop through our rooms, and grab any messages, and write / broadcast them

	}

}

func readFromClient(c *Connection, r *Room, ch chan (string)) {
	// read input from client and do something i dont know what yet
	// defer c.connection.Close()
	for {

		// fmt.Println("reading from user")
		nr := bufio.NewReader(c.connection)

		s := bufio.NewScanner(nr)
		for s.Scan() {
			// r.messages <- s.Text()
			ch <- s.Text()
		}
		// input := s.Text()
		// if len(input) > 0 {
		// 	fmt.Printf("input: %v", input)
		// 	r.messages <- input
		// }
		// fmt.Printf("input %v", len(input))

		// if err != nil {
		// 	fmt.Printf("Error: %v", err)
		// 	break
		// }
		// fmt.Println("read from user")
		// fmt.Println("sending to channel")

		// fmt.Println(input)
	}

}

func writeToClient(c *Connection, msg string, id *string) {
	// Write to our client
	if _, err := io.WriteString(c.connection, "-"+*id+">"+msg+">\n>"); err != nil {
		fmt.Printf("Error: %v", err)
	}
	// done <- true
}

// func readFromClient(c *Connection, r *[]Room) {
// 	// read input from client and do something i dont know what yet
// 	fmt.Println("Reading from client")
// 	for {
// 		input, err := bufio.NewReader(c.connection).ReadString('\n')
// 		if err != nil {
// 			fmt.Printf("Error: %v", err)
// 			break
// 		}

// 		// broadcast user message if in a specific room? lets func that
// 		if c.room == "coolblue" {
// 			fmt.Println("broadcasting in coolblue")
// 			go broadcast(c, &(*r)[0], "-"+input)

// 		}
// 		if c.room == "bottleup" {
// 			fmt.Println("broadcasting in bottleup")
// 			go broadcast(c, &(*r)[0], "-"+input)

// 		}

// 		if c.room == "none" {
// 			sendBackToClient(c.connection, ">you are not in a room so nobody read this...\n>")
// 		}

// 		// fmt.Println(input)
// 	}

// }

func broadcast(c *Connection, r *Room, s string) {

	for i := 0; i < len((*r).users); i++ {

		// lets send some data back over the connection if NOT current user
		fmt.Printf("%v = %v\n", (*r).users[i].id, c.id)
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
