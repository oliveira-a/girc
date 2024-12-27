package main

import (
	"fmt"
	"log"
	"net"
)

const PORT = ":2000"

type Command struct {
	msg string
}

func main() {
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Println("Listening on port 2000.")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handle(conn)
	}
}

func handle(c net.Conn) {
	defer c.Close()

	b := make([]byte, 1024)
	n, err := c.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Message from '%s': %s",
		c.RemoteAddr(), b[:n])
}
