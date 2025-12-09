// internal/tui/keymap/binding.go

package keymap

import "strings"

// Mode represents the editor mode.
type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeVisual
	ModeOperatorPending // Waiting for motion after operator (d, y, c)
)

// String returns the mode name.
func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeVisual:
		return "VISUAL"
	case ModeOperatorPending:
		return "OPERATOR"
	default:
		return "UNKNOWN"
	}
}

// ParseMode converts a string to a Mode.
func ParseMode(s string) Mode {
	switch strings.ToLower(s) {
	case "normal":
		return ModeNormal
	case "insert":
		return ModeInsert
	case "visual":
		return ModeVisual
	case "operator", "operator_pending":
		return ModeOperatorPending
	default:
		return ModeNormal
	}
}

// Binding represents a key-to-action binding.
type Binding struct {
	Key    string // Key or key sequence (e.g., "j", "gg", "ctrl+s")
	Action Action // Action to perform
}

// BindingMap holds all bindings for a specific mode.
type BindingMap struct {
	// Single key bindings (fast lookup)
	single map[string]Action

	// Multi-key sequence bindings (e.g., "gg", "dd")
	// Maps first key to possible sequences
	sequences map[string][]sequenceBinding
}

// sequenceBinding is a multi-key binding.
type sequenceBinding struct {
	keys   string // Full key sequence
	action Action
}

// NewBindingMap creates a new empty binding map.
func NewBindingMap() *BindingMap {
	return &BindingMap{
		single:    make(map[string]Action),
		sequences: make(map[string][]sequenceBinding),
	}
}

// Bind adds a key binding.
func (b *BindingMap) Bind(key string, action Action) {
	if len(key) == 0 {
		return
	}

	if len(key) == 1 || isSpecialKey(key) {
		// Single key or special key (ctrl+x, etc.)
		b.single[key] = action
	} else {
		// Multi-key sequence
		firstKey := string(key[0])
		b.sequences[firstKey] = append(b.sequences[firstKey], sequenceBinding{
			keys:   key,
			action: action,
		})
	}
}

// Unbind removes a key binding.
func (b *BindingMap) Unbind(key string) {
	if len(key) == 1 || isSpecialKey(key) {
		delete(b.single, key)
	} else {
		firstKey := string(key[0])
		seqs := b.sequences[firstKey]
		for i, seq := range seqs {
			if seq.keys == key {
				b.sequences[firstKey] = append(seqs[:i], seqs[i+1:]...)
				break
			}
		}
	}
}

// LookupResult represents the result of looking up a key.
type LookupResult struct {
	// Action is the resolved action (ActionNone if not found or pending)
	Action Action

	// Status indicates the lookup status
	Status LookupStatus

	// PossibleSequences are sequences that could still match
	PossibleSequences []string
}

// LookupStatus indicates the result of a key lookup.
type LookupStatus int

const (
	// LookupNotFound - no binding found
	LookupNotFound LookupStatus = iota

	// LookupFound - exact match found
	LookupFound

	// LookupPending - could be start of a multi-key sequence
	LookupPending

	// LookupPartialMatch - matches part of a sequence, waiting for more
	LookupPartialMatch
)

// Lookup looks up a key or key sequence.
func (b *BindingMap) Lookup(keys string) LookupResult {
	if len(keys) == 0 {
		return LookupResult{Status: LookupNotFound}
	}

	// Check for exact single key match
	if len(keys) == 1 || isSpecialKey(keys) {
		if action, ok := b.single[keys]; ok {
			// Check if this could also be start of a sequence
			firstKey := keys
			if len(keys) > 1 {
				firstKey = keys // special keys are treated as single
			} else {
				firstKey = string(keys[0])
			}

			if seqs, hasSeqs := b.sequences[firstKey]; hasSeqs && len(seqs) > 0 {
				// Could be single key OR start of sequence
				return LookupResult{
					Action:            action,
					Status:            LookupPending,
					PossibleSequences: getSequenceKeys(seqs),
				}
			}

			return LookupResult{
				Action: action,
				Status: LookupFound,
			}
		}

		// Check if this is start of a sequence
		firstKey := keys
		if len(keys) == 1 {
			firstKey = string(keys[0])
		}

		if seqs, ok := b.sequences[firstKey]; ok && len(seqs) > 0 {
			return LookupResult{
				Status:            LookupPending,
				PossibleSequences: getSequenceKeys(seqs),
			}
		}

		return LookupResult{Status: LookupNotFound}
	}

	// Multi-key sequence lookup
	firstKey := string(keys[0])
	seqs, ok := b.sequences[firstKey]
	if !ok {
		return LookupResult{Status: LookupNotFound}
	}

	var exactMatch Action
	var possibleMatches []string

	for _, seq := range seqs {
		if seq.keys == keys {
			// Exact match
			exactMatch = seq.action
		} else if strings.HasPrefix(seq.keys, keys) {
			// Partial match - could continue
			possibleMatches = append(possibleMatches, seq.keys)
		}
	}

	if exactMatch != ActionNone {
		if len(possibleMatches) > 0 {
			// Exact match but could also continue
			return LookupResult{
				Action:            exactMatch,
				Status:            LookupPending,
				PossibleSequences: possibleMatches,
			}
		}
		return LookupResult{
			Action: exactMatch,
			Status: LookupFound,
		}
	}

	if len(possibleMatches) > 0 {
		return LookupResult{
			Status:            LookupPartialMatch,
			PossibleSequences: possibleMatches,
		}
	}

	return LookupResult{Status: LookupNotFound}
}

// GetAction returns the action for a key, or ActionNone.
// This is a simple lookup that doesn't handle sequences.
func (b *BindingMap) GetAction(key string) Action {
	if action, ok := b.single[key]; ok {
		return action
	}
	return ActionNone
}

// AllBindings returns all bindings as a slice.
func (b *BindingMap) AllBindings() []Binding {
	var bindings []Binding

	// Single key bindings
	for key, action := range b.single {
		bindings = append(bindings, Binding{Key: key, Action: action})
	}

	// Sequence bindings
	for _, seqs := range b.sequences {
		for _, seq := range seqs {
			bindings = append(bindings, Binding{Key: seq.keys, Action: seq.action})
		}
	}

	return bindings
}

// Clone creates a copy of the binding map.
func (b *BindingMap) Clone() *BindingMap {
	clone := NewBindingMap()

	for key, action := range b.single {
		clone.single[key] = action
	}

	for key, seqs := range b.sequences {
		clone.sequences[key] = make([]sequenceBinding, len(seqs))
		copy(clone.sequences[key], seqs)
	}

	return clone
}

// isSpecialKey checks if a key is a special key (ctrl+x, alt+x, etc.).
func isSpecialKey(key string) bool {
	return strings.HasPrefix(key, "ctrl+") ||
		strings.HasPrefix(key, "alt+") ||
		strings.HasPrefix(key, "shift+") ||
		key == "enter" ||
		key == "tab" ||
		key == "backspace" ||
		key == "delete" ||
		key == "esc" ||
		key == "up" ||
		key == "down" ||
		key == "left" ||
		key == "right" ||
		key == "home" ||
		key == "end" ||
		key == "pgup" ||
		key == "pgdown" ||
		key == "f1" ||
		key == "f2" ||
		key == "f3" ||
		key == "f4" ||
		key == "f5" ||
		key == "f6" ||
		key == "f7" ||
		key == "f8" ||
		key == "f9" ||
		key == "f10" ||
		key == "f11" ||
		key == "f12"
}

// getSequenceKeys extracts the key strings from sequence bindings.
func getSequenceKeys(seqs []sequenceBinding) []string {
	keys := make([]string, len(seqs))
	for i, seq := range seqs {
		keys[i] = seq.keys
	}
	return keys
}

// NormalizeKey normalizes a key string for consistent lookup.
func NormalizeKey(key string) string {
	// Convert to lowercase for special keys
	lower := strings.ToLower(key)

	// Normalize common variations
	switch lower {
	case "return":
		return "enter"
	case "escape":
		return "esc"
	case "pageup":
		return "pgup"
	case "pagedown":
		return "pgdown"
	case "space", " ":
		return "space"
	}

	// Normalize ctrl/alt/shift combinations
	if strings.Contains(lower, "ctrl+") ||
		strings.Contains(lower, "alt+") ||
		strings.Contains(lower, "shift+") {
		return lower
	}

	// Keep original case for regular keys
	return key
}