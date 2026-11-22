package audit

import (
	"sync"
	"time"
)

// Entry is one audited admin operation.
type Entry struct {
	Time   time.Time `json:"time"`
	User   string    `json:"user"`
	Method string    `json:"method"`
	Path   string    `json:"path"`
	Status int       `json:"status"`
}

// Store is a bounded in-memory audit log.
type Store struct {
	mu      sync.Mutex
	max     int
	entries []Entry
}

// NewStore creates a store keeping at most max entries (default 1000).
func NewStore(max int) *Store {
	if max <= 0 {
		max = 1000
	}
	return &Store{max: max}
}

// Add appends an entry, trimming the oldest when over capacity.
func (s *Store) Add(e Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, e)
	if len(s.entries) > s.max {
		s.entries = s.entries[len(s.entries)-s.max:]
	}
}

// List returns newest-first entries with optional method filter and pagination.
func (s *Store) List(limit, offset int, method string) (items []Entry, total int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	filtered := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		if method != "" && e.Method != method {
			continue
		}
		filtered = append(filtered, e)
	}
	total = len(filtered)
	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	}
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 50
	}
	if offset >= len(filtered) {
		return []Entry{}, total
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total
}

// Clear removes all entries.
func (s *Store) Clear() {
	s.mu.Lock()
	s.entries = nil
	s.mu.Unlock()
}
