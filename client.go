package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	hostFlag  = flag.String("h", "", "the irc host. i.e. localhost")
	portFlag  = flag.String("p", "", "the irc server port number. i.e. 2000.")
	aliasFlag = flag.String("a", "", "an alias. i.e. andreb")

	messages []string
	conn     net.Conn
	mu       sync.Mutex
)

func init() {
	flag.Parse()
}

func main() {
	var err error
	conn, err = net.Dial("tcp", *hostFlag+":"+*portFlag)
	if err != nil {
		log.Fatal(err)
	}

	connect()

	go handleServerInput()

	reader := bufio.NewScanner(os.Stdin)
	c := 0
	for {
		clearScreen()

		printMessages()

		c += 1
		fmt.Print(*aliasFlag + strconv.Itoa(c) + "> ")
		if !reader.Scan() {
			break
		}
		msg := strings.TrimSpace(reader.Text())

		err = sendMessage(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func printMessages() {
	for _, msg := range messages {
		fmt.Println(msg)
	}
}

func connect() {
	conn.Write([]byte("CONNECT " + *aliasFlag))
}

func sendMessage(m string) error {
	_, err := conn.Write([]byte("MESSAGE " + m))

	return err
}

func clearScreen() {
	fmt.Print("\033[H\033[J")
}

func handleServerInput() {
	for {
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		mu.Lock()

		messages = append(messages, string(b[:n]))

		mu.Unlock()
	}
}
