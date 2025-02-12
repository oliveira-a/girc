package main

import (
	"log"
	"net"
	"strings"
	"sync"
)

type Client struct {
	// The alias associated with the client.
	Alias string

	// The connection associated with this client.
	Connection net.Conn
}

// Use this method to message a client through its existing connection.
func (c *Client) Message(from string, content string) {
	c.Connection.Write([]byte(from + ": " + content + "\n"))
}

// Kills the connection to the client.
func (c *Client) Disconnect() {
	c.Connection.Close()
}

const PORT = "2000"

var (
	clientPool = make(map[string]*Client)
	mu         sync.Mutex
)

func main() {
	l, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Printf("Listening on port %s...", PORT)

	for {
		// Wait for incoming connections.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Connection has been received, handle it in its
		// own go-routine.
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	var alias string

	for {
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		msg := strings.Fields(string(b[:n]))
		cmd, content := strings.ToUpper(msg[0]), strings.Join(msg[1:], " ")

		switch cmd {
		case "CONNECT":
			mu.Lock()

			// Clients must have a unique alias.
			// todo: ensure content is an acceptable alias (e.g. no spaces).
			_, exists := clientPool[content]
			if !exists {
				clientPool[content] = &Client{
					Alias:      content,
					Connection: conn,
				}

				log.Printf("Client with alias '%s' added to the client pool.\n", content)

				alias = content
			} else {
				log.Printf("Client with alias '%s' already exists.\n", content)
			}

			mu.Unlock()
			break
		case "MESSAGE":
			if alias != "" {
				for _, c := range clientPool {
					c.Message(alias, content)
				}
			} else {
				conn.Write([]byte("Alias not set. Use `CONNECT <alias>`.\n"))
			}
			break
		default:
			log.Printf("Command type '%s' not supported.\n", cmd)
		}
	}
}
