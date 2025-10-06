package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		currentline := ""
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			parts := strings.Split(string(b), "\n")
			currentline += parts[0]
			if len(parts) == 1 {
			} else {
				currentline = strings.Replace(currentline, "\n", "", 1)
				lines <- currentline
				// fmt.Printf("read: %s\n", currentline)
				currentline = parts[1]
			}
			if err != nil || n == 0 {
				break
			}
		}
		close(lines)
	}()

	return lines
}

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		for line := range getLinesChannel(conn) {
			fmt.Printf("%s\n", line)
		}
	}
}
