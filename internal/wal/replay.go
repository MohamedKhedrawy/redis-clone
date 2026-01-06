package wal

import (
    "bufio"
    "encoding/binary"
    "io"
    "os"
)

// Replay reads WAL records from the file at path and applies them using the provided apply function.

func Replay(path string, apply func([]byte) error) error {
    f, err := os.Open(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return err
    }
    defer f.Close()

    r := bufio.NewReader(f)

    for {
        var lenBuf [4]byte
        if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
            if err == io.EOF || err == io.ErrUnexpectedEOF {
                return nil // safe stop
            }
            return err
        }

        n := binary.BigEndian.Uint32(lenBuf[:])
        rec := make([]byte, n)

        if _, err := io.ReadFull(r, rec); err != nil {
            return nil // partial tail record → ignore
        }

        if err := apply(rec); err != nil {
            return err
        }
    }
}
