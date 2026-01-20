// internal/eval/userfunc.go

package eval

import (
	"strings"
	"sync"
	"time"

	"github.com/0xsj/numio/internal/ast"
	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// USER FUNCTION
// ════════════════════════════════════════════════════════════════

// UserFunction represents a user-defined function at runtime.
type UserFunction struct {
	Name        string    // Function name
	Params      []string  // Parameter names
	Body        ast.Expr  // Body expression (AST)
	BodySource  string    // Original body source (for persistence/display)
	Description string    // Optional user description
	CreatedAt   time.Time // Creation timestamp
	UpdatedAt   time.Time // Last update timestamp
}

// NewUserFunction creates a new user function from a FuncDefStmt.
func NewUserFunction(def *ast.FuncDefStmt) *UserFunction {
	now := time.Now()
	return &UserFunction{
		Name:       def.Name,
		Params:     def.Params,
		Body:       def.Body,
		BodySource: extractBodySource(def),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// NewUserFunctionFromParts creates a user function from individual components.
// Useful for restoring from persistence.
func NewUserFunctionFromParts(name string, params []string, body ast.Expr, bodySource string) *UserFunction {
	now := time.Now()
	return &UserFunction{
		Name:       name,
		Params:     params,
		Body:       body,
		BodySource: bodySource,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// extractBodySource extracts the body source from a FuncDefStmt.
func extractBodySource(def *ast.FuncDefStmt) string {
	if def.Body != nil {
		return def.Body.String()
	}
	return ""
}

// Arity returns the number of parameters.
func (f *UserFunction) Arity() int {
	return len(f.Params)
}

// MinArgs returns the minimum number of arguments (same as Arity for now).
func (f *UserFunction) MinArgs() int {
	return len(f.Params)
}

// MaxArgs returns the maximum number of arguments (same as Arity for now).
func (f *UserFunction) MaxArgs() int {
	return len(f.Params)
}

// String returns a string representation of the function definition.
func (f *UserFunction) String() string {
	var sb strings.Builder
	sb.WriteString("def ")
	sb.WriteString(f.Name)
	sb.WriteString("(")
	sb.WriteString(strings.Join(f.Params, ", "))
	sb.WriteString("): ")
	sb.WriteString(f.BodySource)
	return sb.String()
}

// Signature returns just the function signature (name and params).
func (f *UserFunction) Signature() string {
	var sb strings.Builder
	sb.WriteString(f.Name)
	sb.WriteString("(")
	sb.WriteString(strings.Join(f.Params, ", "))
	sb.WriteString(")")
	return sb.String()
}

// Clone creates a deep copy of the user function.
func (f *UserFunction) Clone() *UserFunction {
	params := make([]string, len(f.Params))
	copy(params, f.Params)

	return &UserFunction{
		Name:        f.Name,
		Params:      params,
		Body:        f.Body, // AST nodes are generally immutable
		BodySource:  f.BodySource,
		Description: f.Description,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}
}

// ════════════════════════════════════════════════════════════════
// USER FUNCTION REGISTRY
// ════════════════════════════════════════════════════════════════

// UserFuncRegistry holds user-defined functions in memory.
type UserFuncRegistry struct {
	mu        sync.RWMutex
	functions map[string]*UserFunction
}

// NewUserFuncRegistry creates a new empty registry.
func NewUserFuncRegistry() *UserFuncRegistry {
	return &UserFuncRegistry{
		functions: make(map[string]*UserFunction),
	}
}

// ════════════════════════════════════════════════════════════════
// REGISTRY OPERATIONS
// ════════════════════════════════════════════════════════════════

// Define registers a new user function or updates an existing one.
// Returns an error message if the function name conflicts with a built-in.
func (r *UserFuncRegistry) Define(fn *UserFunction) string {
	if fn == nil {
		return "function is nil"
	}

	name := strings.ToLower(fn.Name)

	// Check for conflict with built-in functions
	if HasFunction(name) {
		return "cannot redefine built-in function: " + fn.Name
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Update timestamp if replacing
	if existing, ok := r.functions[name]; ok {
		fn.CreatedAt = existing.CreatedAt
		fn.UpdatedAt = time.Now()
	}

	r.functions[name] = fn
	return ""
}

// Get retrieves a user function by name.
func (r *UserFuncRegistry) Get(name string) (*UserFunction, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fn, ok := r.functions[strings.ToLower(name)]
	return fn, ok
}

// Has checks if a user function exists.
func (r *UserFuncRegistry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.functions[strings.ToLower(name)]
	return ok
}

// Delete removes a user function.
// Returns true if the function existed and was removed.
func (r *UserFuncRegistry) Delete(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	name = strings.ToLower(name)
	if _, ok := r.functions[name]; ok {
		delete(r.functions, name)
		return true
	}
	return false
}

// List returns all user function names.
func (r *UserFuncRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.functions))
	for name := range r.functions {
		names = append(names, name)
	}
	return names
}

// All returns all user functions.
func (r *UserFuncRegistry) All() []*UserFunction {
	r.mu.RLock()
	defer r.mu.RUnlock()

	funcs := make([]*UserFunction, 0, len(r.functions))
	for _, fn := range r.functions {
		funcs = append(funcs, fn.Clone())
	}
	return funcs
}

// Count returns the number of user functions.
func (r *UserFuncRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.functions)
}

// Clear removes all user functions.
func (r *UserFuncRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.functions = make(map[string]*UserFunction)
}

// Clone creates a copy of the registry.
func (r *UserFuncRegistry) Clone() *UserFuncRegistry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clone := NewUserFuncRegistry()
	for name, fn := range r.functions {
		clone.functions[name] = fn.Clone()
	}
	return clone
}

// ════════════════════════════════════════════════════════════════
// FUNCTION INVOCATION
// ════════════════════════════════════════════════════════════════

// Call invokes a user function with the given arguments.
// The evalFunc parameter is used to evaluate the function body with bindings.
// This allows the registry to remain decoupled from the evaluator.
func (r *UserFuncRegistry) Call(
	name string,
	args []types.Value,
	evalFunc func(body ast.Expr, bindings map[string]types.Value) types.Value,
) types.Value {
	fn, ok := r.Get(name)
	if !ok {
		return types.Errorf("undefined function: %s", name)
	}

	// Validate argument count
	if len(args) != fn.Arity() {
		if fn.Arity() == 1 {
			return types.Errorf("%s requires exactly 1 argument, got %d", fn.Name, len(args))
		}
		return types.Errorf("%s requires exactly %d arguments, got %d", fn.Name, fn.Arity(), len(args))
	}

	// Build parameter bindings
	bindings := make(map[string]types.Value, len(fn.Params))
	for i, param := range fn.Params {
		bindings[param] = args[i]
	}

	// Evaluate body with bindings
	return evalFunc(fn.Body, bindings)
}

// ════════════════════════════════════════════════════════════════
// SERIALIZATION HELPERS
// ════════════════════════════════════════════════════════════════

// FunctionData represents a user function for serialization.
type FunctionData struct {
	Name        string   `json:"name"`
	Params      []string `json:"params"`
	Body        string   `json:"body"`
	Description string   `json:"description,omitempty"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}

// ToData converts a UserFunction to serializable form.
func (f *UserFunction) ToData() FunctionData {
	return FunctionData{
		Name:        f.Name,
		Params:      f.Params,
		Body:        f.BodySource,
		Description: f.Description,
		CreatedAt:   f.CreatedAt.Unix(),
		UpdatedAt:   f.UpdatedAt.Unix(),
	}
}

// ExportAll exports all functions as serializable data.
func (r *UserFuncRegistry) ExportAll() []FunctionData {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data := make([]FunctionData, 0, len(r.functions))
	for _, fn := range r.functions {
		data = append(data, fn.ToData())
	}
	return data
}
