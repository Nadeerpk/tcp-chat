// client/main.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Read welcome message or name prompt
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Send user input
	input := bufio.NewReader(os.Stdin)
	for {
		text, err := input.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Fprint(conn, text)
	}
}
