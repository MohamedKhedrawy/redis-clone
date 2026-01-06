package store

import (
	"sync"
	"time"
	"strings"
	"log"
	"runtime/debug"
	"context"
)

type Store struct {
	data map[string]Value
	mu   sync.RWMutex
}

type Value struct {
	dataValue string
	expiresAt time.Time
}

func (v *Value) GetValue() string {
	return v.dataValue
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Value),
		mu:   sync.RWMutex{},
	}
}

// Added TTL support in Set and Get methods
func (s *Store) Get(key string) (Value, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.data[key]
	if value.expiresAt.IsZero() == false && time.Now().After(value.expiresAt) {
		// Key has expired
		delete(s.data, key)
		return Value{}, false, nil
	}
	return value, exists, nil
}

func (s *Store) Set(args []string, isReplay bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := args[0]
	value := args[1]
	if len(args) == 4 && strings.ToUpper(args[2]) == "EX" {
		expireSeconds, err := time.ParseDuration(args[3] + "s")
		if err != nil {
			return err
		}
		
		s.data[key] = Value{dataValue: value, expiresAt: time.Now().Add(expireSeconds)}
		return nil
	}
	if len(args) == 4 && strings.ToUpper(args[2]) == "PX" {
		expireMilliseconds, err := time.ParseDuration(args[3] + "ms")
		if err != nil {
			return err
		}
		s.data[key] = Value{dataValue: value, expiresAt: time.Now().Add(expireMilliseconds)}
		return nil
	}
	s.data[key] = Value{dataValue: value, expiresAt: time.Time{}}
	return nil
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}

func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	return keys
}

func (s *Store) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]Value)
	return nil
}

func (s *Store) Exists(key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[key]
	return exists, nil
}

func goSafe(name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic in %s: %v\n%s", name, r, debug.Stack())
			}
		}()
		fn()
	}()
}

func (s *Store) StartCollector(ctx context.Context) {
	goSafe("ttl-expired-collector", func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.collectExpiredKeys()
		}
	})
}

func (s *Store) collectExpiredKeys() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for key, value := range s.data {
		if !value.expiresAt.IsZero() && now.After(value.expiresAt) {
			delete(s.data, key)
		}
	}
}
