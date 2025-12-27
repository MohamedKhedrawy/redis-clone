package parser

import (
	"errors"
	"strings"

	"github.com/MohamedKhedrawy/redis-clone/api/store"
)

type Command struct {
	Name string
	Args []string
}

func ParseCmd(s *store.Store, message string) (string, error) {
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
	value, err := handleCmd(cmd, s)
	if err != nil {
		return "", err
	} else {
		return value, nil
	}
}

func handleCmd(cmd Command, s *store.Store) (string, error) {
	switch cmd.Name {
	case "SET":
		if len(cmd.Args) >= 2 {
			err := handleSet(s, cmd.Args)
			if err != nil {
				return "", err
			} else {
				return "ok", nil
			}
		} else {
			return "", errors.New("SET command requires at least 2 arguments")
		}
	case "GET":
		if len(cmd.Args) == 1 {
			value, err := handleGet(s, cmd.Args)
			if err != nil {
				return "", err
			} else {
				return value, nil
			}
		} else {
			return "", errors.New("GET command requires exactly 1 argument")
		}
	case "DEL":
		if len(cmd.Args) > 0 {
			err := handleDel(s, cmd.Args)
			if err != nil {
				return "", err
			} else {
				return "ok", nil
			}
		} else {
			return "", errors.New("DEL command requires at least 1 argument")
		}
	case "EXISTS":
		if len(cmd.Args) == 1 {
			isExists, err := handleExists(s, cmd.Args)
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

func handleSet(s *store.Store, args []string) error {
	err := s.Set(args)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func handleGet(s *store.Store, args []string) (string, error) {

	value, exists, err := s.Get(args[0])
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New("key does not exist " + args[0])
	}
	return value.GetValue(), nil
}

func handleDel(s *store.Store, args []string) error {
	err := s.Delete(args[0])
	if err != nil {
		return err
	} else {
		return nil
	}
}

func handleExists(s *store.Store, args []string) (bool, error) {
	isExists, err := s.Exists(args[0])
	if err != nil {
		return false, err
	} else {
		return isExists, nil
	}
}
