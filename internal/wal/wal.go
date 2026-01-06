package wal


type WAL interface {
	Append(entry []byte) error
	Close() error
}

