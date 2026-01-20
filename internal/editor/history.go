package editor

// ════════════════════════════════════════════════════════════════
// SNAPSHOT
// ════════════════════════════════════════════════════════════════

// Snapshot represents a point-in-time state of the editor.
type Snapshot struct {
	// Lines is a copy of the buffer content.
	Lines []string

	// CursorRow is the cursor row at this snapshot.
	CursorRow int

	// CursorCol is the cursor column at this snapshot.
	CursorCol int
}

// NewSnapshot creates a snapshot from the current state.
func NewSnapshot(lines []string, cursorRow, cursorCol int) *Snapshot {
	// Deep copy lines
	linesCopy := make([]string, len(lines))
	copy(linesCopy, lines)

	return &Snapshot{
		Lines:     linesCopy,
		CursorRow: cursorRow,
		CursorCol: cursorCol,
	}
}

// Clone creates a copy of the snapshot.
func (s *Snapshot) Clone() *Snapshot {
	linesCopy := make([]string, len(s.Lines))
	copy(linesCopy, s.Lines)

	return &Snapshot{
		Lines:     linesCopy,
		CursorRow: s.CursorRow,
		CursorCol: s.CursorCol,
	}
}

// ════════════════════════════════════════════════════════════════
// HISTORY
// ════════════════════════════════════════════════════════════════

// History manages undo/redo state for the editor.
type History struct {
	undoStack []*Snapshot
	redoStack []*Snapshot

	// maxSize limits the number of undo states kept.
	// 0 means unlimited.
	maxSize int
}

// NewHistory creates a new history manager.
func NewHistory() *History {
	return &History{
		undoStack: nil,
		redoStack: nil,
		maxSize:   0,
	}
}

// NewHistoryWithLimit creates a history with a maximum size.
func NewHistoryWithLimit(maxSize int) *History {
	return &History{
		undoStack: nil,
		redoStack: nil,
		maxSize:   maxSize,
	}
}

// ════════════════════════════════════════════════════════════════
// OPERATIONS
// ════════════════════════════════════════════════════════════════

// Save saves the current state to the undo stack.
// This should be called BEFORE making a change.
func (h *History) Save(lines []string, cursorRow, cursorCol int) {
	snapshot := NewSnapshot(lines, cursorRow, cursorCol)
	h.undoStack = append(h.undoStack, snapshot)

	// Clear redo stack on new change
	h.redoStack = nil

	// Enforce max size
	if h.maxSize > 0 && len(h.undoStack) > h.maxSize {
		// Remove oldest entries
		excess := len(h.undoStack) - h.maxSize
		h.undoStack = h.undoStack[excess:]
	}
}

// SaveSnapshot saves an existing snapshot to the undo stack.
func (h *History) SaveSnapshot(snapshot *Snapshot) {
	h.undoStack = append(h.undoStack, snapshot.Clone())

	// Clear redo stack on new change
	h.redoStack = nil

	// Enforce max size
	if h.maxSize > 0 && len(h.undoStack) > h.maxSize {
		excess := len(h.undoStack) - h.maxSize
		h.undoStack = h.undoStack[excess:]
	}
}

// Undo restores the previous state.
// Returns the snapshot to restore, or nil if no undo available.
// The currentState should be the state BEFORE undoing (to save to redo stack).
func (h *History) Undo(currentLines []string, cursorRow, cursorCol int) *Snapshot {
	if len(h.undoStack) == 0 {
		return nil
	}

	// Save current state to redo stack
	current := NewSnapshot(currentLines, cursorRow, cursorCol)
	h.redoStack = append(h.redoStack, current)

	// Pop from undo stack
	idx := len(h.undoStack) - 1
	snapshot := h.undoStack[idx]
	h.undoStack = h.undoStack[:idx]

	return snapshot
}

// Redo restores a previously undone state.
// Returns the snapshot to restore, or nil if no redo available.
// The currentState should be the state BEFORE redoing (to save to undo stack).
func (h *History) Redo(currentLines []string, cursorRow, cursorCol int) *Snapshot {
	if len(h.redoStack) == 0 {
		return nil
	}

	// Save current state to undo stack
	current := NewSnapshot(currentLines, cursorRow, cursorCol)
	h.undoStack = append(h.undoStack, current)

	// Pop from redo stack
	idx := len(h.redoStack) - 1
	snapshot := h.redoStack[idx]
	h.redoStack = h.redoStack[:idx]

	return snapshot
}

// ════════════════════════════════════════════════════════════════
// QUERIES
// ════════════════════════════════════════════════════════════════

// CanUndo returns true if there are undo states available.
func (h *History) CanUndo() bool {
	return len(h.undoStack) > 0
}

// CanRedo returns true if there are redo states available.
func (h *History) CanRedo() bool {
	return len(h.redoStack) > 0
}

// UndoCount returns the number of undo states available.
func (h *History) UndoCount() int {
	return len(h.undoStack)
}

// RedoCount returns the number of redo states available.
func (h *History) RedoCount() int {
	return len(h.redoStack)
}

// ════════════════════════════════════════════════════════════════
// MANAGEMENT
// ════════════════════════════════════════════════════════════════

// Clear removes all history.
func (h *History) Clear() {
	h.undoStack = nil
	h.redoStack = nil
}

// ClearRedo removes only the redo stack.
func (h *History) ClearRedo() {
	h.redoStack = nil
}

// SetMaxSize sets the maximum undo history size.
// If the current history exceeds this, oldest entries are removed.
func (h *History) SetMaxSize(maxSize int) {
	h.maxSize = maxSize

	if maxSize > 0 && len(h.undoStack) > maxSize {
		excess := len(h.undoStack) - maxSize
		h.undoStack = h.undoStack[excess:]
	}
}

// MaxSize returns the maximum history size (0 = unlimited).
func (h *History) MaxSize() int {
	return h.maxSize
}

// ════════════════════════════════════════════════════════════════
// PEEK
// ════════════════════════════════════════════════════════════════

// PeekUndo returns the most recent undo state without removing it.
// Returns nil if no undo available.
func (h *History) PeekUndo() *Snapshot {
	if len(h.undoStack) == 0 {
		return nil
	}
	return h.undoStack[len(h.undoStack)-1]
}

// PeekRedo returns the most recent redo state without removing it.
// Returns nil if no redo available.
func (h *History) PeekRedo() *Snapshot {
	if len(h.redoStack) == 0 {
		return nil
	}
	return h.redoStack[len(h.redoStack)-1]
}

// ════════════════════════════════════════════════════════════════
// CHANGE GROUP
// ════════════════════════════════════════════════════════════════

// ChangeGroup allows multiple edits to be undone as a single operation.
type ChangeGroup struct {
	history   *History
	startSize int
	committed bool
}

// BeginGroup starts a change group.
// All changes until EndGroup are undone together.
func (h *History) BeginGroup() *ChangeGroup {
	return &ChangeGroup{
		history:   h,
		startSize: len(h.undoStack),
		committed: false,
	}
}

// EndGroup ends the change group, collapsing all changes into one.
func (g *ChangeGroup) EndGroup() {
	if g.committed {
		return
	}
	g.committed = true

	// If we have more than one new state, collapse them
	newStates := len(g.history.undoStack) - g.startSize
	if newStates <= 1 {
		return // Nothing to collapse
	}

	// Keep only the first state of the group
	// (the state before any changes in this group)
	g.history.undoStack = g.history.undoStack[:g.startSize+1]
}

// Cancel cancels the change group, removing all states added during it.
func (g *ChangeGroup) Cancel() {
	if g.committed {
		return
	}
	g.committed = true

	// Remove all states added during this group
	if len(g.history.undoStack) > g.startSize {
		g.history.undoStack = g.history.undoStack[:g.startSize]
	}
}

// ════════════════════════════════════════════════════════════════
// MARK / RESTORE
// ════════════════════════════════════════════════════════════════

// Mark represents a position in the history.
type Mark struct {
	undoSize int
}

// Mark creates a mark at the current history position.
func (h *History) Mark() Mark {
	return Mark{
		undoSize: len(h.undoStack),
	}
}

// RestoreToMark removes all history after the mark.
// Does not change the buffer - only the history.
func (h *History) RestoreToMark(m Mark) {
	if m.undoSize < len(h.undoStack) {
		h.undoStack = h.undoStack[:m.undoSize]
	}
	h.redoStack = nil
}

// HasChangedSince returns true if the history has changed since the mark.
func (h *History) HasChangedSince(m Mark) bool {
	return len(h.undoStack) != m.undoSize
}
