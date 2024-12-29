package main

import (
	"encoding/json"
	"log"
	"net"
	"strings"
	"sync"
)

const PORT = "2000"

type CommandType int

const (
	Connect CommandType = iota
	Message
)

type Command struct {
	CommandType CommandType `json:"commandType"`
	Msg         string      `json:"msg"`
}

var (
	connectionPool = make(map[string]net.Conn)
	mu             sync.Mutex
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
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

	hostName := getHostName(conn.RemoteAddr().String())

	switch cmd.CommandType {
	case Connect:
		if connectionExists(hostName) {
			log.Printf("Client %s already connected...\n", hostName)
			conn.Close()
			break
		}

		addConnection(hostName, conn)

		log.Printf("Client %s connected...\n", hostName)
		break
	default:
		log.Printf("Command '%s' not supported.\n", cmd.CommandType)
		conn.Close()
	}
}

func addConnection(connName string, c net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	connectionPool[connName] = c
}

func connectionExists(connName string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, exists := connectionPool[connName]

	return exists
}

func getHostName(addr string) string {
	i := strings.LastIndex(addr, ":")

	return addr[:i]
}
