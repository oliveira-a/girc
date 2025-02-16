package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"golang.org/x/term"
)

var (
	hostFlag  = flag.String("h", "", "the irc host. i.e. localhost")
	portFlag  = flag.String("p", "", "the irc server port number. i.e. 2000.")
	aliasFlag = flag.String("a", "", "an alias. i.e. andreb")
)

func init() {
	flag.Parse()
}

func main() {
	conn, err := net.Dial("tcp", *hostFlag+":"+*portFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Let the server know of the alias.
	conn.Write([]byte("CONNECT " + *aliasFlag))

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(fd, oldState)

	var wg sync.WaitGroup
	wg.Add(2)

	si := make(chan string)
	ui := make(chan byte)

	// handle incoming messages from the server.
	go func(ch chan string) {
		defer wg.Done()
		for {
			b := make([]byte, 1024)
			n, err := conn.Read(b)
			if err != nil {
				log.Fatal(err)
			}

			s := string(b[:n])
			ch <- strings.TrimRight(s, "\n")
		}
	}(si)

	// handle user input.
	go func(ch chan byte) {
		defer wg.Done()
		b := make([]byte, 1)

		for {
			_, err := os.Stdin.Read(b)
			if err != nil {
				log.Fatal(err)
			}

			ch <- b[0]
		}
	}(ui)

	var messages []string
	var ib []byte

	clearScreen()
	prompt(ib)

	for {
		select {
		case msg := <-si:
			messages = append(messages, msg)

			clearScreen()
			moveCursorToOrigin()

			for _, m := range messages {
				fmt.Print(m)
				moveCursorNLinesDown(1)
			}

			prompt(ib)
		case input := <-ui:
			switch input {
            case 27:
                return
			case 127:
				if len(ib) > 0 {
					ib = ib[:len(ib)-1]
                    clearLine()
                    prompt(ib)
				}
			case 13:
				if len(ib) > 0 {
					_, err := conn.Write([]byte("MESSAGE " + string(ib)))
					if err != nil {
						log.Fatal(err)
					}
					ib = nil
                    clearLine()
					moveCursorToLineBegin()
					prompt(ib)
				}
			default:
				ib = append(ib, input)
				fmt.Print(string(input))
			}
		}
	}

	wg.Wait()
}

func prompt(i []byte) {
    // move the cursor to the bottom left of the screen.
	printEscape("999E")

	fmt.Printf("%s> %s", *aliasFlag, string(i))
}

func moveCursorNLinesDown(n int) {
	printEscape(string(n) + "E")
}

func moveCursorToOrigin() {
	printEscape("H")
}

func moveCursorToLineBegin() {
    printEscape("999D")
}
func clearScreen() {
	printEscape("2J")
}

func clearLine() {
	printEscape("2K")
}

func printEscape(ec string) {
	fmt.Print("\033[" + ec)
}

