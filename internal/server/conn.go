package server

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	parser "github.com/MohamedKhedrawy/redis-clone/internal/command"
	"github.com/MohamedKhedrawy/redis-clone/internal/store"
	"github.com/MohamedKhedrawy/redis-clone/internal/wal"
)

func handleConnection(conn net.Conn, KvStore *store.Store, w wal.WAL) {
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

		// Set read deadline
		conn.SetReadDeadline(time.Now().Add(readTimeout))

		// Read the length prefix (4 bytes)
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

		// Set write deadline
		conn.SetWriteDeadline(time.Now().Add(writeTimeout))

		// Process the message using the parser
		response, err := executeWithTimeout(10*time.Second, func() (string, error) {
			return parser.ParseCmd(KvStore, w, message, false)
		})

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
