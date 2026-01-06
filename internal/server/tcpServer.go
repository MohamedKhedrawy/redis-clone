package server

import (
	"fmt"
	"net"
	"github.com/MohamedKhedrawy/redis-clone/internal/store"
	"github.com/MohamedKhedrawy/redis-clone/internal/wal"
	parser "github.com/MohamedKhedrawy/redis-clone/internal/command"
	"context"
	"time"
	"errors"
)

// var mut sync.RWMutex
const MaxMessageSize = 1048576 // 1 MB
const (
	readTimeout  = 2 * time.Minute
	writeTimeout = 10 * time.Second
)

func RunServer() {
	KvStore := store.NewStore()
	w, err := wal.OpenFileWAL("wal.log")
	if err != nil {
		fmt.Println("Error opening WAL:", err)
		return
	}
	defer w.Close()

	err = wal.Replay("wal.log", func(rec []byte) error {
		_, err := parser.ParseCmd(KvStore, w, string(rec), true)
		return err
	})
	if err != nil {
		fmt.Println("Error replaying WAL:", err)
		return
	}

	KvStore.StartCollector(context.Background())
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
		go handleConnection(conn, KvStore, w)
	}

}

func executeWithTimeout(
	timeout time.Duration,
	fn func() (string, error),
) (string, error) {

	done := make(chan struct{})
	var res string
	var err error

	go func() {
		defer close(done)
		res, err = fn()
	}()

	select {
	case <-done:
		return res, err
	case <-time.After(timeout):
		return "", errors.New("command timeout")
	}
}
