package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

func listen() (net.Listener, error) {
	if os.Getenv("LISTEN_PID") == strconv.Itoa(os.Getpid()) && os.Getenv("LISTEN_FDS") == "1" {
		f := os.NewFile(3, "systemd-socket")
		defer f.Close()
		return net.FileListener(f)
	}

	log.Fatal("Expected systemd socket..")
	return nil, os.ErrPermission
}

func main() {
	connection, err := listen()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening..")

	conn, err := connection.Accept()
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			return
		}
		log.Printf("Error: %v", err)
		os.Exit(3)
	}
	handle(conn)
}

func handle(conn net.Conn) {
	defer conn.Close()

	log.Println("Got new connection")
	if _, err := io.Copy(conn, conn); err != nil {
		log.Printf("Something went wrong: %v", err)
	}
	log.Println("End of connection")
}
