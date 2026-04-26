package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	connection, err := net.Dial("tcp", "localhost:4000")
	buf := make([]byte, 1024)

	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	connection.Write([]byte("hello "))
	connection.Write([]byte("world"))

	n, netErr := connection.Read(buf)
	if netErr != nil {
		log.Fatal(netErr)
	}

	fmt.Printf("got %d bytes: %q\n", n, buf[:n])

}
