package editor

// ════════════════════════════════════════════════════════════════
// CURSOR
// ════════════════════════════════════════════════════════════════

// Cursor represents the cursor position in the editor.
type Cursor struct {
	row int
	col int

	// Preferred column for vertical movement.
	// When moving up/down, the cursor tries to maintain this column.
	preferredCol int
}

// NewCursor creates a cursor at position (0, 0).
func NewCursor() *Cursor {
	return &Cursor{
		row:          0,
		col:          0,
		preferredCol: 0,
	}
}

// NewCursorAt creates a cursor at the specified position.
func NewCursorAt(row, col int) *Cursor {
	if row < 0 {
		row = 0
	}
	if col < 0 {
		col = 0
	}
	return &Cursor{
		row:          row,
		col:          col,
		preferredCol: col,
	}
}

// ════════════════════════════════════════════════════════════════
// ACCESSORS
// ════════════════════════════════════════════════════════════════

// Row returns the current row (0-indexed).
func (c *Cursor) Row() int {
	return c.row
}

// Col returns the current column (0-indexed).
func (c *Cursor) Col() int {
	return c.col
}

// Position returns the current (row, col) position.
func (c *Cursor) Position() (int, int) {
	return c.row, c.col
}

// PreferredCol returns the preferred column for vertical movement.
func (c *Cursor) PreferredCol() int {
	return c.preferredCol
}

// ════════════════════════════════════════════════════════════════
// SETTERS
// ════════════════════════════════════════════════════════════════

// SetPosition sets the cursor position.
func (c *Cursor) SetPosition(row, col int) {
	if row < 0 {
		row = 0
	}
	if col < 0 {
		col = 0
	}
	c.row = row
	c.col = col
	c.preferredCol = col
}

// SetRow sets the row, preserving preferred column behavior.
func (c *Cursor) SetRow(row int) {
	if row < 0 {
		row = 0
	}
	c.row = row
}

// SetCol sets the column and updates preferred column.
func (c *Cursor) SetCol(col int) {
	if col < 0 {
		col = 0
	}
	c.col = col
	c.preferredCol = col
}

// SetColPreservePreferred sets the column without updating preferred column.
// Used for vertical movement.
func (c *Cursor) SetColPreservePreferred(col int) {
	if col < 0 {
		col = 0
	}
	c.col = col
}

// ════════════════════════════════════════════════════════════════
// MOVEMENT
// ════════════════════════════════════════════════════════════════

// MoveUp moves the cursor up by n rows.
// Returns true if the cursor moved.
func (c *Cursor) MoveUp(n int) bool {
	if c.row == 0 || n <= 0 {
		return false
	}

	c.row -= n
	if c.row < 0 {
		c.row = 0
	}

	return true
}

// MoveDown moves the cursor down by n rows.
// The caller must clamp to buffer bounds.
// Returns true if the cursor moved.
func (c *Cursor) MoveDown(n int) bool {
	if n <= 0 {
		return false
	}

	c.row += n
	return true
}

// MoveLeft moves the cursor left by n columns.
// Returns true if the cursor moved.
func (c *Cursor) MoveLeft(n int) bool {
	if c.col == 0 || n <= 0 {
		return false
	}

	c.col -= n
	if c.col < 0 {
		c.col = 0
	}
	c.preferredCol = c.col

	return true
}

// MoveRight moves the cursor right by n columns.
// The caller must clamp to line bounds.
// Returns true if the cursor moved.
func (c *Cursor) MoveRight(n int) bool {
	if n <= 0 {
		return false
	}

	c.col += n
	c.preferredCol = c.col

	return true
}

// MoveToLineStart moves the cursor to the start of the line.
func (c *Cursor) MoveToLineStart() {
	c.col = 0
	c.preferredCol = 0
}

// MoveToLineEnd moves the cursor to the specified column (end of line).
func (c *Cursor) MoveToLineEnd(lineLength int) {
	c.col = lineLength
	c.preferredCol = lineLength
}

// MoveToFirstNonBlank moves the cursor to the first non-blank character.
func (c *Cursor) MoveToFirstNonBlank(line string) {
	col := 0
	for col < len(line) && isWhitespace(line[col]) {
		col++
	}
	c.col = col
	c.preferredCol = col
}

// ════════════════════════════════════════════════════════════════
// WORD MOVEMENT
// ════════════════════════════════════════════════════════════════

// MoveWordForward moves the cursor to the start of the next word.
// Returns the new column position.
func (c *Cursor) MoveWordForward(line string) int {
	col := c.col

	// Skip current word
	for col < len(line) && isWordChar(line[col]) {
		col++
	}

	// Skip non-word characters (whitespace, punctuation)
	for col < len(line) && !isWordChar(line[col]) {
		col++
	}

	c.col = col
	c.preferredCol = col

	return col
}

// MoveWordBackward moves the cursor to the start of the previous word.
// Returns the new column position.
func (c *Cursor) MoveWordBackward(line string) int {
	col := c.col

	// Move back one if at word boundary
	if col > 0 {
		col--
	}

	// Skip non-word characters backwards
	for col > 0 && !isWordChar(line[col]) {
		col--
	}

	// Skip word characters backwards
	for col > 0 && isWordChar(line[col-1]) {
		col--
	}

	c.col = col
	c.preferredCol = col

	return col
}

// MoveWordEnd moves the cursor to the end of the current/next word.
// Returns the new column position.
func (c *Cursor) MoveWordEnd(line string) int {
	col := c.col

	// Move forward one if at word end
	if col < len(line) {
		col++
	}

	// Skip non-word characters
	for col < len(line) && !isWordChar(line[col]) {
		col++
	}

	// Move to end of word
	for col < len(line) && isWordChar(line[col]) {
		col++
	}

	// Back up one to be ON the last character
	if col > 0 && col <= len(line) {
		col--
	}

	c.col = col
	c.preferredCol = col

	return col
}

// ════════════════════════════════════════════════════════════════
// CLAMPING
// ════════════════════════════════════════════════════════════════

// Clamp ensures the cursor is within valid bounds for the buffer.
// maxRow is the last valid row index.
// lineLength is the length of the current line.
// allowPastEnd allows cursor to be at lineLength (insert mode).
func (c *Cursor) Clamp(maxRow, lineLength int, allowPastEnd bool) {
	// Clamp row
	if c.row < 0 {
		c.row = 0
	}
	if c.row > maxRow {
		c.row = maxRow
	}

	// Clamp column
	maxCol := lineLength
	if !allowPastEnd && maxCol > 0 {
		maxCol-- // Normal mode: cursor on last char, not past it
	}

	if c.col < 0 {
		c.col = 0
	}
	if c.col > maxCol {
		c.col = maxCol
	}
	if c.col < 0 {
		c.col = 0 // Handle empty lines
	}
}

// ClampCol clamps only the column to the given line length.
func (c *Cursor) ClampCol(lineLength int, allowPastEnd bool) {
	maxCol := lineLength
	if !allowPastEnd && maxCol > 0 {
		maxCol--
	}

	if c.col > maxCol {
		c.col = maxCol
	}
	if c.col < 0 {
		c.col = 0
	}
}

// ApplyPreferredCol applies the preferred column, clamped to line length.
func (c *Cursor) ApplyPreferredCol(lineLength int, allowPastEnd bool) {
	maxCol := lineLength
	if !allowPastEnd && maxCol > 0 {
		maxCol--
	}

	if c.preferredCol <= maxCol {
		c.col = c.preferredCol
	} else {
		c.col = maxCol
	}
	if c.col < 0 {
		c.col = 0
	}
}

// ════════════════════════════════════════════════════════════════
// UTILITY
// ════════════════════════════════════════════════════════════════

// Clone creates a copy of the cursor.
func (c *Cursor) Clone() *Cursor {
	return &Cursor{
		row:          c.row,
		col:          c.col,
		preferredCol: c.preferredCol,
	}
}

// Equals returns true if two cursors have the same position.
func (c *Cursor) Equals(other *Cursor) bool {
	if other == nil {
		return false
	}
	return c.row == other.row && c.col == other.col
}

// IsBefore returns true if this cursor is before the other cursor.
func (c *Cursor) IsBefore(other *Cursor) bool {
	if other == nil {
		return false
	}
	if c.row < other.row {
		return true
	}
	if c.row > other.row {
		return false
	}
	return c.col < other.col
}

// IsAfter returns true if this cursor is after the other cursor.
func (c *Cursor) IsAfter(other *Cursor) bool {
	if other == nil {
		return false
	}
	if c.row > other.row {
		return true
	}
	if c.row < other.row {
		return false
	}
	return c.col > other.col
}

// ════════════════════════════════════════════════════════════════
// FIND OPERATIONS
// ════════════════════════════════════════════════════════════════

// FindCharForward finds the next occurrence of a character on the current line.
// Returns the column if found, -1 otherwise.
func (c *Cursor) FindCharForward(line string, ch byte) int {
	for i := c.col + 1; i < len(line); i++ {
		if line[i] == ch {
			return i
		}
	}
	return -1
}

// FindCharBackward finds the previous occurrence of a character on the current line.
// Returns the column if found, -1 otherwise.
func (c *Cursor) FindCharBackward(line string, ch byte) int {
	for i := c.col - 1; i >= 0; i-- {
		if line[i] == ch {
			return i
		}
	}
	return -1
}

// MoveToChar moves to the specified column if valid.
// Returns true if moved.
func (c *Cursor) MoveToChar(col int) bool {
	if col < 0 {
		return false
	}
	c.col = col
	c.preferredCol = col
	return true
}
