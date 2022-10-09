package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

func handleError(err error) {
	if err == io.EOF {
		panic("server closed unexpectedly")
	} else if err != nil {
		panic(err)
	}
}

func read(conn net.Conn) {
	//In a continuous loop, read a message from the server and display it.
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		handleError(err)
		fmt.Print(msg)
	}
}

func write(conn net.Conn) {
	//Continually get input from the user and send messages to the server.
	stdin := bufio.NewReader(os.Stdin)
	for {
		msg, err := stdin.ReadString('\n')
		handleError(err)
		fmt.Fprintf(conn, msg)
	}
}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()
	conn, err := net.Dial("tcp", *addrPtr)
	handleError(err)
	end := make(chan bool)

	go read(conn)
	go write(conn)
	<-end
}
