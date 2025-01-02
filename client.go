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

    r := bufio.NewReader(os.Stdin)
    for {
        fmt.Print(*aliasFlag + "> ")
        m, _ := r.ReadString('\n')

        fmt.Fprintf(c, m)
    }
}
