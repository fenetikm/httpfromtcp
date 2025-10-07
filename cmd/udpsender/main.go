package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udp, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Couldn't resolve UDP address.")
	}

	conn, err := net.DialUDP("udp", nil, udp)
	if err != nil {
		log.Fatalf("Couldn't create UDP connection.")
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Fatalf("Error reading from Stdin")
		}
		_, err = conn.Write(line)
		if err != nil {
			fmt.Printf("err %v\n", err)
		}
	}
}
