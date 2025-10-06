package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Error opening file.")
	}
	defer file.Close()

	currentline := ""
	for {
		b := make([]byte, 8)
		n, err := file.Read(b)
		parts := strings.Split(string(b), "\n")
		currentline += parts[0]
		if len(parts) == 1 {
		} else {
			currentline = strings.Replace(currentline, "\n", "", 1)
			fmt.Printf("read: %s\n", currentline)
			currentline = parts[1]
		}
		if err != nil || n == 0 {
			break
		}
	}
}
