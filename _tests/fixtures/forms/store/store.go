package store

import "sync"

var (
	mu      sync.Mutex
	entries []string
)

func Add(entry string) {
	mu.Lock()
	defer mu.Unlock()
	entries = append(entries, entry)
}

func All() []string {
	mu.Lock()
	defer mu.Unlock()
	return append([]string(nil), entries...)
}
