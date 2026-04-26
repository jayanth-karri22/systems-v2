package main

import (
	"bufio"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:4000")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	conn.Write([]byte("Hello \n"))
	conn.Write([]byte("World \n"))

	for i := 0; i < 2; i++ {
		if !scanner.Scan() {
			break
		}
		log.Printf("got: %q", scanner.Text())
	}

}
