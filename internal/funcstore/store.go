// pkg/funcstore/store.go

// Package funcstore provides interfaces and implementations for
// persisting user-defined functions.
package funcstore

import (
	"errors"
	"time"
)

// ════════════════════════════════════════════════════════════════
// ERRORS
// ════════════════════════════════════════════════════════════════

var (
	// ErrNotFound is returned when a function is not found.
	ErrNotFound = errors.New("function not found")

	// ErrAlreadyExists is returned when trying to create a function that exists.
	ErrAlreadyExists = errors.New("function already exists")

	// ErrInvalidName is returned when a function name is invalid.
	ErrInvalidName = errors.New("invalid function name")

	// ErrStoreClosed is returned when operating on a closed store.
	ErrStoreClosed = errors.New("store is closed")
)

// ════════════════════════════════════════════════════════════════
// STORED FUNCTION
// ════════════════════════════════════════════════════════════════

// StoredFunction represents a user-defined function in storage.
// This is the persistence-layer representation, decoupled from AST.
type StoredFunction struct {
	// Name is the function identifier (unique key).
	Name string `json:"name"`

	// Params is the list of parameter names.
	Params []string `json:"params"`

	// Body is the function body as source code.
	// This will be re-parsed when loaded.
	Body string `json:"body"`

	// Description is an optional user-provided description.
	Description string `json:"description,omitempty"`

	// CreatedAt is when the function was first defined.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when the function was last modified.
	UpdatedAt time.Time `json:"updated_at"`
}

// NewStoredFunction creates a new StoredFunction with timestamps set.
func NewStoredFunction(name string, params []string, body string) *StoredFunction {
	now := time.Now()
	return &StoredFunction{
		Name:      name,
		Params:    params,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Clone creates a deep copy of the stored function.
func (f *StoredFunction) Clone() *StoredFunction {
	params := make([]string, len(f.Params))
	copy(params, f.Params)

	return &StoredFunction{
		Name:        f.Name,
		Params:      params,
		Body:        f.Body,
		Description: f.Description,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}
}

// Signature returns the function signature string.
func (f *StoredFunction) Signature() string {
	sig := f.Name + "("
	for i, p := range f.Params {
		if i > 0 {
			sig += ", "
		}
		sig += p
	}
	sig += ")"
	return sig
}

// Definition returns the full function definition string.
func (f *StoredFunction) Definition() string {
	return "def " + f.Signature() + ": " + f.Body
}

// ════════════════════════════════════════════════════════════════
// FUNCTION STORE INTERFACE
// ════════════════════════════════════════════════════════════════

// FunctionStore defines the interface for function persistence.
// Implementations can use SQLite, key-value stores, files, etc.
type FunctionStore interface {
	// Get retrieves a function by name.
	// Returns ErrNotFound if the function doesn't exist.
	Get(name string) (*StoredFunction, error)

	// Set stores a function, creating or updating as needed.
	Set(fn *StoredFunction) error

	// Delete removes a function by name.
	// Returns ErrNotFound if the function doesn't exist.
	Delete(name string) error

	// List returns all stored function names.
	List() ([]string, error)

	// All returns all stored functions.
	All() ([]*StoredFunction, error)

	// Has checks if a function exists.
	Has(name string) (bool, error)

	// Count returns the number of stored functions.
	Count() (int, error)

	// Clear removes all stored functions.
	Clear() error

	// Close releases any resources held by the store.
	Close() error
}

// ════════════════════════════════════════════════════════════════
// OPTIONAL INTERFACES
// ════════════════════════════════════════════════════════════════

// FunctionStoreWithSearch extends FunctionStore with search capability.
type FunctionStoreWithSearch interface {
	FunctionStore

	// Search finds functions matching a query string.
	// The query may match against name, description, or body.
	Search(query string) ([]*StoredFunction, error)
}

// FunctionStoreWithExport extends FunctionStore with import/export.
type FunctionStoreWithExport interface {
	FunctionStore

	// Export returns all functions as a JSON byte slice.
	Export() ([]byte, error)

	// Import loads functions from a JSON byte slice.
	// If merge is true, existing functions are preserved.
	// If merge is false, the store is cleared first.
	Import(data []byte, merge bool) error
}

// FunctionStoreWithWatch extends FunctionStore with change notifications.
type FunctionStoreWithWatch interface {
	FunctionStore

	// Watch returns a channel that receives function names when they change.
	// The channel is closed when the store is closed.
	Watch() <-chan string
}

// ════════════════════════════════════════════════════════════════
// STORE OPTIONS
// ════════════════════════════════════════════════════════════════

// StoreOptions configures a function store.
type StoreOptions struct {
	// Path is the file path for file-based stores.
	Path string

	// AutoSave enables automatic saving after modifications.
	AutoSave bool

	// SyncWrites forces synchronous writes (slower but safer).
	SyncWrites bool
}

// DefaultStoreOptions returns sensible defaults.
func DefaultStoreOptions() StoreOptions {
	return StoreOptions{
		AutoSave:   true,
		SyncWrites: false,
	}
}
