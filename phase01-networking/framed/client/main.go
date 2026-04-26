package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	for _, msg := range []string{"Hello", "World"} {
		if err := sendFramed(conn, msg); err != nil {
			log.Fatalf("send %q: %v", msg, err)
		}
		reply, err := recvFramed(conn)
		if err != nil {
			log.Fatalf("recv: %v", err)
		}
		log.Printf("got: %q", reply)
	}
}

func sendFramed(conn net.Conn, msg string) error {
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(msg)))
	if _, err := conn.Write(lenBuf); err != nil {
		return err
	}
	if _, err := conn.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}

func recvFramed(conn net.Conn) ([]byte, error) {
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenBuf)
	if length > 1<<20 {
		return nil, io.ErrShortBuffer // any sentinel; we'll do typed errors in Phase 3
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, err
	}
	return payload, nil
}
