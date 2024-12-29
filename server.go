package main

import (
	"encoding/json"
	"log"
	"net"
	"sync"
)

const PORT = "2000"

type CommandType int

const (
	Connect CommandType = iota
	Message
)

type Command struct {
	CommandType CommandType `json:"type"`
	From        string      `json:"from"`
	Message     string      `json:"message"`
}

type Client struct {
	// The nickname the client wants to be associated with.
	Alias string

	// The connection associated with this client.
	Connection net.Conn
}

// Use this method to message a client through its existing connection.
func (c *Client) Message(msg string) {
	c.Connection.Write([]byte(msg))
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

	var cmd Command
	err = json.Unmarshal(b[:n], &cmd)
	if err != nil {
		log.Fatal(err)
	}

	switch cmd.CommandType {
	case Connect:
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
	case Message:
		for _, c := range clientPool {
			log.Printf("Message from '%s': %s\n", cmd.From, cmd.Message)
			c.Message(cmd.Message)
		}
		conn.Close()
		break
	default:
		log.Printf("Command type '%s' not supported.\n", cmd.CommandType)
		conn.Close()
	}
}
