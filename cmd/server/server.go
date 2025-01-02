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
func (c *Client) Message(msg string) {
	m := &shared.ClientMessage{
		From:    c.Alias,
		Content: msg,
	}

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
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
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
		defer mu.Unlock()

		// Ensure the client with the same alias does not already exist
		// since they must be unique in our client pool.
		_, exists := clientPool[cmd.From]
		if exists {
			log.Printf("Client with alias '%s' already exists.\n", cmd.From)
			conn.Close()
			break
		}

		clientPool[cmd.From] = &Client{
			Alias:      cmd.From,
			Connection: conn,
		}

		log.Printf("Client with alias '%s' added to the client pool.\n", cmd.From)
		break
	case shared.Message:
		log.Printf("Message from '%s': %s\n", cmd.From, cmd.Content)
		for _, c := range clientPool {
			c.Message(cmd.Content)
		}
		break
	default:
		log.Printf("Command type '%s' not supported.\n", cmd.CommandType)
		conn.Close()
	}
}
