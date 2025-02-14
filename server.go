package main

import (
	"errors"
	"io"
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
		conn, err := l.Accept()
		log.Printf("Connection acknowledged from client '%s'.", conn.RemoteAddr().String())
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
	addr := conn.RemoteAddr().String()

	for {
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			// client has disconnected or exited unexpectedly.
			if err == io.EOF {
				if _, ok := clientPool[addr]; ok {
					delete(clientPool, addr)
				}
				log.Printf("Client with address '%s' has disconnected.", addr)
				return
			}
		}

		msg := strings.Fields(string(b[:n]))
		cmd, content := strings.ToUpper(msg[0]), strings.Join(msg[1:], " ")

		switch cmd {
		case "CONNECT":
			mu.Lock()

			if err := validateAlias(content); err == nil {
				clientPool[addr] = &Client{
					Alias:      content,
					Connection: conn,
				}
				alias = content
			} else {
				conn.Write([]byte(err.Error()))
				return
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

func validateAlias(a string) error {
	if a == "" {
		return errors.New("Alias cannot be empty.")
	}

	for _, client := range clientPool {
		if client.Alias == a {
			return errors.New("The specified alias already exists.")
		}
	}

	return nil
}
