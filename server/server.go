package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// Deal with an error event.
	if err != nil {
		panic(err)
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for {
		conn, err := ln.Accept()
		handleError(err)
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	reader := bufio.NewReader(client)
	for {
		in, err := reader.ReadString('\n')
		if err != nil {
			client.Close()
			break
		}
		text := fmt.Sprintf("[%d]: %s", clientid, in)
		msg := Message{sender: clientid, message: text}
		msgs <- msg
	}

}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//Create a Listener for TCP connections on the port given above.
	ln, err := net.Listen("tcp", *portPtr)
	handleError(err)

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)
	//Create integer to track current ID to use
	currentID := 0

	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			clients[currentID] = conn
			go handleClient(conn, currentID, msgs)
			currentID++

		case msg := <-msgs:
			for clientID, client := range clients {
				if msg.sender != clientID {
					fmt.Fprintf(client, msg.message)
				}
			}
		}
	}
}
