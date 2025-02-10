package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

    "github.com/pterm/pterm"

	"github.com/oliveira-a/girc/internal/shared"
)

var hostFlag = flag.String("h", "", "the irc host. i.e. localhost")
var portFlag = flag.String("p", "", "the irc server port number. i.e. 2000.")
var aliasFlag = flag.String("a", "", "an alias. i.e. andreb")

func init() {
	flag.Parse()
}

func main() {
	c, err := net.Dial("tcp", *hostFlag+":"+*portFlag)
	if err != nil {
		log.Fatal(err)
	}

	go handleServerInput(c)

	// todo: refactor this in to a helper function
	cmd := &shared.Command{
		CommandType: shared.Connect,
		From:        *aliasFlag,
		Content:     "",
	}
	b, err := json.Marshal(cmd)
	c.Write(b)

	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(*aliasFlag + "> ")
		m, _ := r.ReadString('\n')

		cmd := &shared.Command{
			CommandType: shared.Message,
			From:        *aliasFlag,
			Content:     m,
		}
		b, err := json.Marshal(cmd)
		if err != nil {
			log.Fatal(err)
		}
		_, err = c.Write(b)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func handleServerInput(c net.Conn) {
	for {
		b := make([]byte, 1024)
		n, err := c.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		var msg shared.ClientMessage
		err = json.Unmarshal(b[:n], &msg)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s: %s", msg.From, msg.Content)
	}
}
