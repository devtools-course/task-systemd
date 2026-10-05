package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
)

func main() {
	socketPath := "/var/run/reckod.sock"

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	server := bufio.NewReader(conn)
	message := fmt.Sprintf("Hi, %d", rand.Int())

	if _, err := fmt.Fprintln(conn, message); err != nil {
		log.Fatalf("Could not send: %v", err)
	}

	reply, err := server.ReadString('\n')
	if err != nil {
		log.Fatalf("Could not receive: %v", err)
	}

	if strings.Trim(reply, "\n") != message {
		log.Fatalf("Expected <%q>\nGot <%q>", message, reply)
	}
}
