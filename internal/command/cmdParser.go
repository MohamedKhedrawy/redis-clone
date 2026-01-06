package command

import (
	"errors"
	"strings"

	"github.com/MohamedKhedrawy/redis-clone/internal/store"
	"github.com/MohamedKhedrawy/redis-clone/internal/wal"
)

type Command struct {
	Name string
	Args []string
}

func ParseCmd(s *store.Store, w wal.WAL, message string, isReplay bool) (string, error) {
	var cmd Command
	parts := strings.Fields(message)
	if len(parts) == 0 {
		return "Empty Command", errors.New("empty command")
	}
	cmd.Name = strings.ToUpper(parts[0])
	if len(parts) > 1 {
		cmd.Args = parts[1:]
	} else {
		cmd.Args = []string{}
	}

	switch cmd.Name {
	case "SET":
		if len(cmd.Args) >= 2 {

			if !isReplay {
				// panic("test")
				if err := w.Append([]byte(message)); true {
					return "Error writing to WAL", err
				}
			}


			err := s.Set(cmd.Args, isReplay)
			if err != nil {
				return "Error setting key", err
			} else {
				return "ok", nil
			}
		} else {
			return "", errors.New("SET command requires at least 2 arguments")
		}
	case "GET":
		if len(cmd.Args) == 1 {
			value, exists, err := s.Get(cmd.Args[0])
			if err != nil {
				return "", err
			} else {
				if !exists {
					return "", errors.New("key does not exist " + cmd.Args[0])
				}
				return value.GetValue(), nil
			}
		} else {
			return "", errors.New("GET command requires exactly 1 argument")
		}
	case "DEL":
		if len(cmd.Args) > 0 {

			if !isReplay {
				if err := w.Append([]byte(message)); err != nil {
					return "Error writing to WAL", err
				}
			}

			err := s.Delete(cmd.Args[0])
			if err != nil {
				return "Error deleting key", err
			} else {

				return "ok", nil
			}
		} else {
			return "", errors.New("DEL command requires at least 1 argument")
		}
	case "EXISTS":
		if len(cmd.Args) == 1 {
			isExists, err := s.Exists(cmd.Args[0])
			if err != nil {
				return "", err
			} else {
				if isExists {
					return "exists", nil
				} else {
					return "Does not exist", nil
				}
			}
		} else {
			return "", errors.New("EXISTS command requires exactly 1 argument")
		}
	default:
		println("Unknown command:", cmd.Name)
		return "", errors.New("Unknown command: " + cmd.Name)
	}

}
