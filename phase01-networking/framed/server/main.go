package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

const maxMessageSize = 1 << 20 // 1 MiB

func main() {
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	log.Printf("listening on %s", listener.Addr())

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
	defer conn.Close()
	remote := conn.RemoteAddr()
	log.Printf("client connected: %s", remote)
	defer log.Printf("client disconnected: %s", remote)

	lenBuf := make([]byte, 4)
	for {
		if _, err := io.ReadFull(conn, lenBuf); err != nil {
			if !errors.Is(err, io.EOF) {
				log.Printf("read length from %s: %v", remote, err)
			}
			return
		}

		length := binary.BigEndian.Uint32(lenBuf)
		if length > maxMessageSize {
			log.Printf("message from %s too large: %d bytes", remote, length)
			return
		}

		payload := make([]byte, length)
		if _, err := io.ReadFull(conn, payload); err != nil {
			log.Printf("read payload from %s: %v", remote, err)
			return
		}

		log.Printf("received %d bytes from %s: %q", length, remote, payload)

		if _, err := conn.Write(lenBuf); err != nil {
			log.Printf("write length to %s: %v", remote, err)
			return
		}
		if _, err := conn.Write(payload); err != nil {
			log.Printf("write payload to %s: %v", remote, err)
			return
		}
	}
}
