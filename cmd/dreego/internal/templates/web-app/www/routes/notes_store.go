package routes

import "sync"

var (
	notesMu sync.Mutex
	notes   []string
)

func addNote(title string) {
	notesMu.Lock()
	defer notesMu.Unlock()
	notes = append(notes, title)
}

func notesSnapshot() []string {
	notesMu.Lock()
	defer notesMu.Unlock()
	return append([]string(nil), notes...)
}
