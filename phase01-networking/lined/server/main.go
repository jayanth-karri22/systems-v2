package main

import (
	"bufio"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":4000")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()

		conn.Write([]byte(line + "\n"))
	}
}
