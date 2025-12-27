package main

import (
    "encoding/binary"
    "fmt"
    "log"
    "net"
    "time"
    "bufio"
    "os"
    "io"
)

const (
    SERVER_ADDR = "localhost:8080" // Change to your server address
    CONN_TYPE   = "tcp"
)

func ConstructMessage(message string) []byte {
    body := []byte(message)
    bodyLen := uint32(len(body))

    lengthPrefix := make([]byte, 4)
    binary.BigEndian.PutUint32(lengthPrefix, bodyLen)

    return append(lengthPrefix, body...)
}

func SendSlowly(conn net.Conn, data []byte) error {
    chunkSize := 8
    fmt.Printf("\n--- Sending %d bytes in chunks ---\n", len(data))

    for i := 0; i < len(data); i += chunkSize {
        end := i + chunkSize
        if end > len(data) {
            end = len(data)
        }
        chunk := data[i:end]
        n, err := conn.Write(chunk)
        if err != nil {
            return fmt.Errorf("failed to write chunk: %w", err)
        }
        fmt.Printf("Sent %d bytes... ", n)
        time.Sleep(500 * time.Millisecond) // Faster delay for testing
    }
    fmt.Println("\nMessage sent.")
    return nil
}

func SendMsg(message string) {
    fmt.Println("Connecting to:", SERVER_ADDR)
    conn, err := net.Dial(CONN_TYPE, SERVER_ADDR)
    if err != nil {
        log.Fatal("Connection failed:", err)
    }
    defer conn.Close()

    
    // Construct and send the message
    data := ConstructMessage(message)
    if err := SendSlowly(conn, data); err != nil {
        log.Fatal("Error sending message:", err)
    }
    
    // Read the response
    reader := bufio.NewReader(conn)
    
    msgLenBuf := make([]byte, 4)
    if _, err := io.ReadFull(reader, msgLenBuf); err != nil {
        log.Fatal("Error reading response length:", err)
    }
    
    msgLen := binary.BigEndian.Uint32(msgLenBuf)
    respBuf := make([]byte, msgLen)
    if _, err := io.ReadFull(reader, respBuf); err != nil {
        log.Fatal("Error reading response body:", err)
    }
    fmt.Println(string(respBuf))
}


func main() {
    for {
        scanner := bufio.NewScanner(os.Stdin)
        fmt.Print("Enter message to send: ")
        if scanner.Scan() {
            message := scanner.Text()
        SendMsg(message)
        } else {
            fmt.Println("Error reading input:", scanner.Err())
            break
        }
    }
}