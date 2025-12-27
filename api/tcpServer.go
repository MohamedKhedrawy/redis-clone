package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	// "sync"
	// "github.com/MohamedKhedrawy/redis-clone/api/client"
	"github.com/MohamedKhedrawy/redis-clone/api/parser"
	"github.com/MohamedKhedrawy/redis-clone/api/store"
	"context"
)

// var mut sync.RWMutex
const MaxMessageSize = 1048576 // 1 MB

func main() {
	kvStore := store.NewStore()
	kvStore.StartCollector(context.Background())
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP listener:", err)
		return
	}
	defer ln.Close()
	fmt.Println("TCP listener started on :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, kvStore)
	}

}

func handleConnection(conn net.Conn, kvStore *store.Store) {
	defer conn.Close()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	reader := bufio.NewReader(conn)
	fmt.Println("New connection from", conn.RemoteAddr())

	for {
		// Handle the connection (read/write data) here
		limBuf := make([]byte, 4)

		if _, err := io.ReadFull(reader, limBuf); err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client:", conn.RemoteAddr())
				return
			} else {
				fmt.Println("Error reading from connection:", err)
				return
			}

		}

		length := binary.BigEndian.Uint32(limBuf)

		if length <= 0 {
			fmt.Println("Invalid message length", length)
			return
		}

		if length > MaxMessageSize {
			length = MaxMessageSize
		}

		fmt.Println("Expected message length:", length)

		bodyBuf := make([]byte, length)
		if _, err := io.ReadFull(reader, bodyBuf); err != nil {
			fmt.Println("Error reading message body:", err)
			return
		}

		message := string(bodyBuf)
		message = strings.TrimSpace(message)
		fmt.Println("Received message:", message)

		response, err := parser.ParseCmd(kvStore, message)

		// Echo back the received message
		// response := "Echo: " + message
		if err != nil {
			response = "Error: " + err.Error()
			fmt.Println(response)
		} else {
			fmt.Println(response)
		}
		
		respLen := uint32(len(response))
		respLenBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(respLenBuf, respLen)

		if _, err := conn.Write(respLenBuf); err != nil {
			fmt.Println("Error writing response length:", err)
			return
		}
		
		if _, err := conn.Write([]byte(response)); err != nil {
			fmt.Println("Error writing response body:", err)
			return
		}
	}
}

