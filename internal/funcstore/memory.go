// pkg/funcstore/memory.go

package funcstore

import (
	"encoding/json"
	"sort"
	"strings"
	"sync"
)

// ════════════════════════════════════════════════════════════════
// MEMORY STORE
// ════════════════════════════════════════════════════════════════

// MemoryStore is an in-memory implementation of FunctionStore.
// It is thread-safe and suitable for use during a session.
// Data is lost when the process exits unless explicitly exported.
type MemoryStore struct {
	mu       sync.RWMutex
	funcs    map[string]*StoredFunction
	closed   bool
	watchers []chan string
}

// NewMemoryStore creates a new in-memory function store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		funcs:    make(map[string]*StoredFunction),
		watchers: make([]chan string, 0),
	}
}

// ════════════════════════════════════════════════════════════════
// CORE OPERATIONS
// ════════════════════════════════════════════════════════════════

// Get retrieves a function by name.
func (s *MemoryStore) Get(name string) (*StoredFunction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, ErrStoreClosed
	}

	fn, ok := s.funcs[strings.ToLower(name)]
	if !ok {
		return nil, ErrNotFound
	}

	return fn.Clone(), nil
}

// Set stores a function, creating or updating as needed.
func (s *MemoryStore) Set(fn *StoredFunction) error {
	if fn == nil {
		return ErrInvalidName
	}

	if fn.Name == "" {
		return ErrInvalidName
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStoreClosed
	}

	key := strings.ToLower(fn.Name)

	// Preserve creation time if updating
	if existing, ok := s.funcs[key]; ok {
		fn.CreatedAt = existing.CreatedAt
	}

	s.funcs[key] = fn.Clone()

	// Notify watchers
	s.notifyWatchers(fn.Name)

	return nil
}

// Delete removes a function by name.
func (s *MemoryStore) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStoreClosed
	}

	key := strings.ToLower(name)
	if _, ok := s.funcs[key]; !ok {
		return ErrNotFound
	}

	delete(s.funcs, key)

	// Notify watchers
	s.notifyWatchers(name)

	return nil
}

// List returns all stored function names (sorted alphabetically).
func (s *MemoryStore) List() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, ErrStoreClosed
	}

	names := make([]string, 0, len(s.funcs))
	for _, fn := range s.funcs {
		names = append(names, fn.Name)
	}

	sort.Strings(names)
	return names, nil
}

// All returns all stored functions (sorted by name).
func (s *MemoryStore) All() ([]*StoredFunction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, ErrStoreClosed
	}

	funcs := make([]*StoredFunction, 0, len(s.funcs))
	for _, fn := range s.funcs {
		funcs = append(funcs, fn.Clone())
	}

	// Sort by name
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].Name < funcs[j].Name
	})

	return funcs, nil
}

// Has checks if a function exists.
func (s *MemoryStore) Has(name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return false, ErrStoreClosed
	}

	_, ok := s.funcs[strings.ToLower(name)]
	return ok, nil
}

// Count returns the number of stored functions.
func (s *MemoryStore) Count() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, ErrStoreClosed
	}

	return len(s.funcs), nil
}

// Clear removes all stored functions.
func (s *MemoryStore) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStoreClosed
	}

	s.funcs = make(map[string]*StoredFunction)

	// Notify watchers with empty string to indicate clear
	s.notifyWatchers("")

	return nil
}

// Close releases resources and marks the store as closed.
func (s *MemoryStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStoreClosed
	}

	s.closed = true

	// Close all watcher channels
	for _, ch := range s.watchers {
		close(ch)
	}
	s.watchers = nil

	return nil
}

// ════════════════════════════════════════════════════════════════
// SEARCH (FunctionStoreWithSearch)
// ════════════════════════════════════════════════════════════════

// Search finds functions matching a query string.
// Matches against name, description, and body (case-insensitive).
func (s *MemoryStore) Search(query string) ([]*StoredFunction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, ErrStoreClosed
	}

	if query == "" {
		return s.All()
	}

	query = strings.ToLower(query)
	results := make([]*StoredFunction, 0)

	for _, fn := range s.funcs {
		if strings.Contains(strings.ToLower(fn.Name), query) ||
			strings.Contains(strings.ToLower(fn.Description), query) ||
			strings.Contains(strings.ToLower(fn.Body), query) {
			results = append(results, fn.Clone())
		}
	}

	// Sort by name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

// ════════════════════════════════════════════════════════════════
// EXPORT/IMPORT (FunctionStoreWithExport)
// ════════════════════════════════════════════════════════════════

// Export returns all functions as a JSON byte slice.
func (s *MemoryStore) Export() ([]byte, error) {
	funcs, err := s.All()
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(funcs, "", "  ")
}

// Import loads functions from a JSON byte slice.
// If merge is true, existing functions are preserved (updates allowed).
// If merge is false, the store is cleared first.
func (s *MemoryStore) Import(data []byte, merge bool) error {
	var funcs []*StoredFunction
	if err := json.Unmarshal(data, &funcs); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStoreClosed
	}

	if !merge {
		s.funcs = make(map[string]*StoredFunction)
	}

	for _, fn := range funcs {
		if fn.Name == "" {
			continue
		}
		key := strings.ToLower(fn.Name)

		// Preserve creation time if updating in merge mode
		if merge {
			if existing, ok := s.funcs[key]; ok {
				fn.CreatedAt = existing.CreatedAt
			}
		}

		s.funcs[key] = fn.Clone()
	}

	return nil
}

// ════════════════════════════════════════════════════════════════
// WATCH (FunctionStoreWithWatch)
// ════════════════════════════════════════════════════════════════

// Watch returns a channel that receives function names when they change.
// An empty string indicates the store was cleared.
// The channel is closed when the store is closed.
func (s *MemoryStore) Watch() <-chan string {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan string, 16)
	if !s.closed {
		s.watchers = append(s.watchers, ch)
	} else {
		close(ch)
	}

	return ch
}

// notifyWatchers sends a notification to all watchers.
// Must be called with lock held.
func (s *MemoryStore) notifyWatchers(name string) {
	for _, ch := range s.watchers {
		select {
		case ch <- name:
		default:
			// Channel full, skip (non-blocking)
		}
	}
}

// ════════════════════════════════════════════════════════════════
// INTERFACE COMPLIANCE
// ════════════════════════════════════════════════════════════════

// Compile-time interface checks
var (
	_ FunctionStore           = (*MemoryStore)(nil)
	_ FunctionStoreWithSearch = (*MemoryStore)(nil)
	_ FunctionStoreWithExport = (*MemoryStore)(nil)
	_ FunctionStoreWithWatch  = (*MemoryStore)(nil)
)
