package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
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
	c.Write([]byte("CONNECT " + *aliasFlag))

	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(*aliasFlag + "> ")
		m, _ := r.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}
		_, err = c.Write([]byte("MESSAGE " + m))
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

		fmt.Printf(string(b[:n]))
	}
}
