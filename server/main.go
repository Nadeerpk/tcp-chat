// server/main.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn net.Conn
	name string
}

var (
	clients   = make(map[net.Conn]Client)
	broadcast = make(chan string)
	mutex     sync.Mutex
)

func main() {
	listener, _ := net.Listen("tcp", ":9000")
	defer listener.Close()

	go handleBroadcast()

	for {
		conn, _ := listener.Accept()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("Enter your name: \n"))
	name, _ := bufio.NewReader(conn).ReadString('\n')
	name = strings.TrimSpace(name)

	client := Client{conn: conn, name: name}
	mutex.Lock()
	clients[conn] = client
	mutex.Unlock()

	broadcast <- fmt.Sprintf("%s joined the chat!", name)

	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		msg = strings.TrimSpace(msg)
		if msg != "" {
			broadcast <- fmt.Sprintf("[%s]: %s", name, msg)
		}
	}

	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
	broadcast <- fmt.Sprintf("%s left the chat.", name)
}

func handleBroadcast() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for _, client := range clients {
			client.conn.Write([]byte(msg + "\n"))
		}
		mutex.Unlock()
	}
}
