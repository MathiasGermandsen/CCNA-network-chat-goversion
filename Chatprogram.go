package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn net.Conn
	name string
}

var clients = make(map[client]bool)

func main() {
	fmt.Println("Starting server...")
	listener, err := net.Listen("tcp", "localhost:13000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	fmt.Fprint(conn, "Please enter your name: ")
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading name:", err)
		return
	}
	name = strings.TrimSpace(name)

	cli := client{conn, name}
	clients[cli] = true

	fmt.Printf("%s has joined the chat\n", name)

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		if len(strings.TrimSpace(message)) > 0 {
			fmt.Printf("Received: %s: %s", name, message)
			broadcast(cli, name+": "+message)
		}
	}
}

func broadcast(sender client, message string) {
	for cli := range clients {
		if cli != sender {
			_, err := fmt.Fprint(cli.conn, message)
			if err != nil {
				fmt.Println("Error broadcasting message:", err)
			}
		}
	}
}
