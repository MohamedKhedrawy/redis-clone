package wal

import (
	"bufio"
	"os"
	"sync"
	"encoding/binary"
)


type FileWAL struct {
	mutex    sync.Mutex
	file     *os.File
	buf    *bufio.Writer
}

func OpenFileWAL(filePath string) (WAL, error) {
	file, err := os.OpenFile(filePath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
        0600,)
	if err != nil {
		return nil, err
	}
	return &FileWAL{
		mutex: sync.Mutex{},
		file: file,
		buf:  bufio.NewWriter(file),
	}, nil
}

func (w *FileWAL) Append(entry []byte) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(entry)))

	if _, err := w.buf.Write(lenBuf[:]); err != nil {
		return err
	}
	if _, err := w.buf.Write(entry); err != nil {
		return err
	}
	w.buf.Flush()
	return w.file.Sync()
}

func (w *FileWAL) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	err := w.buf.Flush()
	if err != nil {
		return err
	}
	return w.file.Close()
}
