package main

import (
	"encoding/json"
	"log"
	"net"
	"sync"

	"github.com/oliveira-a/girc/internal/shared"
)

const PORT = "2000"

type Client struct {
	// The alias associated with the client.
	Alias string

	// The connection associated with this client.
	Connection net.Conn
}

// Use this method to message a client through its existing connection.
func (c *Client) Message(m shared.ClientMessage) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	c.Connection.Write(b)
}

// Kills the connection to the client.
func (c *Client) Disconnect() {
	c.Connection.Close()
}

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
	for {
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		var cmd shared.Command
		err = json.Unmarshal(b[:n], &cmd)
		if err != nil {
			log.Fatal(err)
		}

		switch cmd.CommandType {
		case shared.Connect:
			mu.Lock()
			// Clients must have a unique alias.
			_, exists := clientPool[cmd.From]
			if exists {
				log.Printf("Client with alias '%s' already exists.\n", cmd.From)
				conn.Close()
                mu.Unlock()
				break
			}

			clientPool[cmd.From] = &Client{
				Alias:      cmd.From,
				Connection: conn,
			}

			log.Printf("Client with alias '%s' added to the client pool.\n", cmd.From)
            mu.Unlock()
			break
		case shared.Message:
			log.Printf("Message from '%s': %s\n", cmd.From, cmd.Content)
			for _, c := range clientPool {
				cm := &shared.ClientMessage{
					From:    cmd.From,
					Content: cmd.Content,
				}
				c.Message(*cm)
			}
			break
		default:
			log.Printf("Command type '%s' not supported.\n", cmd.CommandType)
			conn.Close()
		}
	}
}
