package main

import (
	"fmt"
	"io"
	"log"
	"os"
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
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Error opening file.")
	}
	defer file.Close()

	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}
