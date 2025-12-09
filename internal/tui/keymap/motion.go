// internal/tui/keymap/motion.go

package keymap

// MotionState tracks the current state of a motion command being built.
// This handles vim-style count + operator + motion sequences like:
//   - "5j"   → count=5, action=MoveDown
//   - "d5j"  → operator=Delete, count=5, motion=MoveDown
//   - "dd"   → operator=Delete, motion=Line (special case)
//   - "y3w"  → operator=Yank, count=3, motion=WordNext
type MotionState struct {
	// Count is the numeric prefix (e.g., "5" in "5j")
	// Zero means no count specified (defaults to 1)
	Count int

	// PendingOperator is set when an operator key is pressed (d, y, c)
	// and we're waiting for a motion
	PendingOperator Action

	// KeyBuffer holds keys for multi-key sequences (e.g., "gg", "dd")
	KeyBuffer string

	// LastAction stores the last completed action for repeat with "."
	LastAction Action
	LastCount  int
	LastMotion Action
}

// NewMotionState creates a new empty motion state.
func NewMotionState() *MotionState {
	return &MotionState{
		Count:           0,
		PendingOperator: ActionNone,
		KeyBuffer:       "",
		LastAction:      ActionNone,
		LastCount:       0,
		LastMotion:      ActionNone,
	}
}

// Reset clears the current motion state.
func (m *MotionState) Reset() {
	m.Count = 0
	m.PendingOperator = ActionNone
	m.KeyBuffer = ""
}

// AddDigit adds a digit to the count.
// Returns true if the digit was consumed.
func (m *MotionState) AddDigit(d rune) bool {
	if d < '0' || d > '9' {
		return false
	}

	// '0' at start is not a count (it's goto line start)
	if d == '0' && m.Count == 0 && m.KeyBuffer == "" {
		return false
	}

	digit := int(d - '0')
	m.Count = m.Count*10 + digit
	return true
}

// GetCount returns the effective count (minimum 1).
func (m *MotionState) GetCount() int {
	if m.Count == 0 {
		return 1
	}
	return m.Count
}

// HasCount returns true if a count was explicitly specified.
func (m *MotionState) HasCount() bool {
	return m.Count > 0
}

// SetOperator sets the pending operator.
func (m *MotionState) SetOperator(op Action) {
	m.PendingOperator = op
}

// HasPendingOperator returns true if there's a pending operator.
func (m *MotionState) HasPendingOperator() bool {
	return m.PendingOperator != ActionNone
}

// AddKey adds a key to the buffer for multi-key sequences.
func (m *MotionState) AddKey(key string) {
	m.KeyBuffer += key
}

// GetKeyBuffer returns the current key buffer.
func (m *MotionState) GetKeyBuffer() string {
	return m.KeyBuffer
}

// ClearKeyBuffer clears the key buffer.
func (m *MotionState) ClearKeyBuffer() {
	m.KeyBuffer = ""
}

// SaveForRepeat saves the current action for repeat with ".".
func (m *MotionState) SaveForRepeat(action Action, count int, motion Action) {
	if action.IsRepeatable() {
		m.LastAction = action
		m.LastCount = count
		m.LastMotion = motion
	}
}

// GetRepeat returns the last repeatable action.
func (m *MotionState) GetRepeat() (Action, int, Action) {
	return m.LastAction, m.LastCount, m.LastMotion
}

// Command represents a fully parsed command ready for execution.
type Command struct {
	// Action is the main action to perform
	Action Action

	// Count is how many times to repeat (1 if not specified)
	Count int

	// Motion is the motion for operator commands (e.g., "w" in "dw")
	Motion Action

	// MotionCount is the count for the motion (e.g., "3" in "d3w")
	MotionCount int

	// Char is for commands that take a character argument (e.g., "f", "t", "r")
	Char rune

	// Raw is the original key sequence
	Raw string
}

// NewCommand creates a simple command with no motion.
func NewCommand(action Action, count int) Command {
	if count < 1 {
		count = 1
	}
	return Command{
		Action: action,
		Count:  count,
	}
}

// NewOperatorCommand creates a command with an operator and motion.
func NewOperatorCommand(operator Action, operatorCount int, motion Action, motionCount int) Command {
	if operatorCount < 1 {
		operatorCount = 1
	}
	if motionCount < 1 {
		motionCount = 1
	}
	return Command{
		Action:      operator,
		Count:       operatorCount,
		Motion:      motion,
		MotionCount: motionCount,
	}
}

// TotalCount returns the total count (count * motion count).
// For "2d3j", this would be 2 * 3 = 6.
func (c Command) TotalCount() int {
	total := c.Count
	if c.MotionCount > 1 {
		total *= c.MotionCount
	}
	if total < 1 {
		total = 1
	}
	return total
}

// IsOperatorPending returns true if this is a partial operator command.
func (c Command) IsOperatorPending() bool {
	return c.Action.IsOperator() && c.Motion == ActionNone
}

// String returns a string representation of the command.
func (c Command) String() string {
	if c.Action == ActionNone {
		return "<none>"
	}

	result := ""

	if c.Count > 1 {
		result += string(rune('0' + c.Count))
	}

	result += string(c.Action)

	if c.Motion != ActionNone {
		if c.MotionCount > 1 {
			result += string(rune('0' + c.MotionCount))
		}
		result += string(c.Motion)
	}

	return result
}

// PendingDisplay returns what to show in the UI for pending commands.
// e.g., "d" when waiting for motion, "5" when count entered
func (m *MotionState) PendingDisplay() string {
	display := ""

	if m.Count > 0 {
		display += intToString(m.Count)
	}

	if m.PendingOperator != ActionNone {
		switch m.PendingOperator {
		case ActionOperatorDelete:
			display += "d"
		case ActionOperatorYank:
			display += "y"
		case ActionOperatorChange:
			display += "c"
		}
	}

	display += m.KeyBuffer

	return display
}

// intToString converts an int to string without fmt package.
func intToString(n int) string {
	if n == 0 {
		return "0"
	}

	var digits []rune
	for n > 0 {
		digits = append([]rune{rune('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}