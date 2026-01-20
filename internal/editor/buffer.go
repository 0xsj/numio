// Package editor provides a headless editor implementation for numio.
// It manages text state, cursor, selection, and produces render output
// that can be consumed by any UI frontend (TUI, native, web).
package editor

// ════════════════════════════════════════════════════════════════
// BUFFER
// ════════════════════════════════════════════════════════════════

// Buffer represents the text content of the editor.
// It stores lines of text and provides methods for manipulation.
type Buffer struct {
	lines []string
}

// NewBuffer creates an empty buffer with one empty line.
func NewBuffer() *Buffer {
	return &Buffer{
		lines: []string{""},
	}
}

// NewBufferFromLines creates a buffer from existing lines.
func NewBufferFromLines(lines []string) *Buffer {
	if len(lines) == 0 {
		return NewBuffer()
	}

	// Copy to avoid external mutation
	copied := make([]string, len(lines))
	copy(copied, lines)

	return &Buffer{
		lines: copied,
	}
}

// NewBufferFromString creates a buffer from a string, splitting on newlines.
func NewBufferFromString(content string) *Buffer {
	if content == "" {
		return NewBuffer()
	}

	lines := splitLines(content)
	return &Buffer{
		lines: lines,
	}
}

// ════════════════════════════════════════════════════════════════
// ACCESSORS
// ════════════════════════════════════════════════════════════════

// LineCount returns the number of lines in the buffer.
func (b *Buffer) LineCount() int {
	return len(b.lines)
}

// Line returns the content of a specific line (0-indexed).
// Returns empty string if row is out of bounds.
func (b *Buffer) Line(row int) string {
	if row < 0 || row >= len(b.lines) {
		return ""
	}
	return b.lines[row]
}

// LineLength returns the length of a specific line.
func (b *Buffer) LineLength(row int) int {
	return len(b.Line(row))
}

// Lines returns a copy of all lines.
func (b *Buffer) Lines() []string {
	copied := make([]string, len(b.lines))
	copy(copied, b.lines)
	return copied
}

// String returns the entire buffer content as a single string.
func (b *Buffer) String() string {
	return joinLines(b.lines)
}

// IsEmpty returns true if the buffer has no content.
func (b *Buffer) IsEmpty() bool {
	return len(b.lines) == 0 || (len(b.lines) == 1 && b.lines[0] == "")
}

// LastRow returns the index of the last row.
func (b *Buffer) LastRow() int {
	if len(b.lines) == 0 {
		return 0
	}
	return len(b.lines) - 1
}

// ════════════════════════════════════════════════════════════════
// CHARACTER OPERATIONS
// ════════════════════════════════════════════════════════════════

// InsertChar inserts a character at the specified position.
func (b *Buffer) InsertChar(row, col int, ch rune) {
	if row < 0 || row >= len(b.lines) {
		return
	}

	line := b.lines[row]

	// Clamp col to valid range
	if col < 0 {
		col = 0
	}
	if col > len(line) {
		col = len(line)
	}

	b.lines[row] = line[:col] + string(ch) + line[col:]
}

// InsertString inserts a string at the specified position.
func (b *Buffer) InsertString(row, col int, s string) {
	if row < 0 || row >= len(b.lines) {
		return
	}

	line := b.lines[row]

	// Clamp col to valid range
	if col < 0 {
		col = 0
	}
	if col > len(line) {
		col = len(line)
	}

	b.lines[row] = line[:col] + s + line[col:]
}

// DeleteChar deletes the character at the specified position.
// Returns the deleted character, or 0 if nothing was deleted.
func (b *Buffer) DeleteChar(row, col int) rune {
	if row < 0 || row >= len(b.lines) {
		return 0
	}

	line := b.lines[row]

	if col < 0 || col >= len(line) {
		return 0
	}

	ch := rune(line[col])
	b.lines[row] = line[:col] + line[col+1:]

	return ch
}

// DeleteCharBack deletes the character before the specified position (backspace).
// Returns the deleted character, or 0 if nothing was deleted.
func (b *Buffer) DeleteCharBack(row, col int) rune {
	if col > 0 {
		return b.DeleteChar(row, col-1)
	}
	return 0
}

// ════════════════════════════════════════════════════════════════
// LINE OPERATIONS
// ════════════════════════════════════════════════════════════════

// SetLine replaces the content of a specific line.
func (b *Buffer) SetLine(row int, content string) {
	if row < 0 || row >= len(b.lines) {
		return
	}
	b.lines[row] = content
}

// InsertLine inserts a new line at the specified row.
// Existing lines are shifted down.
func (b *Buffer) InsertLine(row int, content string) {
	if row < 0 {
		row = 0
	}
	if row > len(b.lines) {
		row = len(b.lines)
	}

	b.lines = append(b.lines[:row], append([]string{content}, b.lines[row:]...)...)
}

// DeleteLine removes the line at the specified row.
// Returns the deleted line content.
func (b *Buffer) DeleteLine(row int) string {
	if row < 0 || row >= len(b.lines) {
		return ""
	}

	deleted := b.lines[row]
	b.lines = append(b.lines[:row], b.lines[row+1:]...)

	// Ensure at least one empty line
	if len(b.lines) == 0 {
		b.lines = []string{""}
	}

	return deleted
}

// SplitLine splits a line at the specified column.
// Content after the column moves to a new line below.
func (b *Buffer) SplitLine(row, col int) {
	if row < 0 || row >= len(b.lines) {
		return
	}

	line := b.lines[row]

	// Clamp col
	if col < 0 {
		col = 0
	}
	if col > len(line) {
		col = len(line)
	}

	before := line[:col]
	after := line[col:]

	b.lines[row] = before
	b.InsertLine(row+1, after)
}

// JoinLines joins the specified line with the line below it.
// Returns the column position where the join occurred.
func (b *Buffer) JoinLines(row int) int {
	if row < 0 || row >= len(b.lines)-1 {
		return 0
	}

	joinCol := len(b.lines[row])
	b.lines[row] = b.lines[row] + b.lines[row+1]
	b.lines = append(b.lines[:row+1], b.lines[row+2:]...)

	return joinCol
}

// AppendLine adds a new line at the end of the buffer.
func (b *Buffer) AppendLine(content string) {
	b.lines = append(b.lines, content)
}

// ════════════════════════════════════════════════════════════════
// RANGE OPERATIONS
// ════════════════════════════════════════════════════════════════

// DeleteRange deletes text from (startRow, startCol) to (endRow, endCol).
// The range is inclusive of start and exclusive of end.
func (b *Buffer) DeleteRange(startRow, startCol, endRow, endCol int) string {
	// Normalize range
	if startRow > endRow || (startRow == endRow && startCol > endCol) {
		startRow, endRow = endRow, startRow
		startCol, endCol = endCol, startCol
	}

	// Bounds check
	if startRow < 0 {
		startRow = 0
	}
	if endRow >= len(b.lines) {
		endRow = len(b.lines) - 1
	}

	// Same line deletion
	if startRow == endRow {
		line := b.lines[startRow]
		if startCol < 0 {
			startCol = 0
		}
		if endCol > len(line) {
			endCol = len(line)
		}
		deleted := line[startCol:endCol]
		b.lines[startRow] = line[:startCol] + line[endCol:]
		return deleted
	}

	// Multi-line deletion
	var deleted string

	// Capture deleted text
	firstLine := b.lines[startRow]
	if startCol > len(firstLine) {
		startCol = len(firstLine)
	}
	deleted = firstLine[startCol:]

	for i := startRow + 1; i < endRow; i++ {
		deleted += "\n" + b.lines[i]
	}

	lastLine := b.lines[endRow]
	if endCol > len(lastLine) {
		endCol = len(lastLine)
	}
	deleted += "\n" + lastLine[:endCol]

	// Perform deletion
	b.lines[startRow] = firstLine[:startCol] + lastLine[endCol:]
	b.lines = append(b.lines[:startRow+1], b.lines[endRow+1:]...)

	return deleted
}

// GetRange returns text from (startRow, startCol) to (endRow, endCol).
func (b *Buffer) GetRange(startRow, startCol, endRow, endCol int) string {
	// Normalize range
	if startRow > endRow || (startRow == endRow && startCol > endCol) {
		startRow, endRow = endRow, startRow
		startCol, endCol = endCol, startCol
	}

	// Bounds check
	if startRow < 0 || startRow >= len(b.lines) {
		return ""
	}
	if endRow >= len(b.lines) {
		endRow = len(b.lines) - 1
	}

	// Same line
	if startRow == endRow {
		line := b.lines[startRow]
		if startCol < 0 {
			startCol = 0
		}
		if endCol > len(line) {
			endCol = len(line)
		}
		if startCol >= endCol {
			return ""
		}
		return line[startCol:endCol]
	}

	// Multi-line
	var result string

	firstLine := b.lines[startRow]
	if startCol > len(firstLine) {
		startCol = len(firstLine)
	}
	result = firstLine[startCol:]

	for i := startRow + 1; i < endRow; i++ {
		result += "\n" + b.lines[i]
	}

	lastLine := b.lines[endRow]
	if endCol > len(lastLine) {
		endCol = len(lastLine)
	}
	result += "\n" + lastLine[:endCol]

	return result
}

// ════════════════════════════════════════════════════════════════
// WORD OPERATIONS
// ════════════════════════════════════════════════════════════════

// DeleteWord deletes from the position to the end of the current word.
// Returns the deleted text and the new column position.
func (b *Buffer) DeleteWord(row, col int) (string, int) {
	if row < 0 || row >= len(b.lines) {
		return "", col
	}

	line := b.lines[row]
	if col >= len(line) {
		return "", col
	}

	endCol := col

	// Skip current word characters
	for endCol < len(line) && isWordChar(line[endCol]) {
		endCol++
	}

	// Skip trailing whitespace
	for endCol < len(line) && isWhitespace(line[endCol]) {
		endCol++
	}

	if endCol == col {
		// Delete at least one character
		endCol = col + 1
		if endCol > len(line) {
			endCol = len(line)
		}
	}

	deleted := line[col:endCol]
	b.lines[row] = line[:col] + line[endCol:]

	return deleted, col
}

// DeleteToLineEnd deletes from the position to the end of the line.
// Returns the deleted text.
func (b *Buffer) DeleteToLineEnd(row, col int) string {
	if row < 0 || row >= len(b.lines) {
		return ""
	}

	line := b.lines[row]
	if col >= len(line) {
		return ""
	}

	deleted := line[col:]
	b.lines[row] = line[:col]

	return deleted
}

// ════════════════════════════════════════════════════════════════
// BULK OPERATIONS
// ════════════════════════════════════════════════════════════════

// Clear resets the buffer to a single empty line.
func (b *Buffer) Clear() {
	b.lines = []string{""}
}

// Replace replaces the entire buffer content.
func (b *Buffer) Replace(content string) {
	if content == "" {
		b.Clear()
		return
	}
	b.lines = splitLines(content)
}

// ReplaceLines replaces all lines.
func (b *Buffer) ReplaceLines(lines []string) {
	if len(lines) == 0 {
		b.Clear()
		return
	}

	b.lines = make([]string, len(lines))
	copy(b.lines, lines)
}

// Clone creates a deep copy of the buffer.
func (b *Buffer) Clone() *Buffer {
	return NewBufferFromLines(b.lines)
}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

// splitLines splits a string into lines.
func splitLines(s string) []string {
	var lines []string
	var current []byte

	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, string(current))
			current = current[:0]
		} else if s[i] == '\r' {
			// Skip \r, handle \r\n as single newline
			if i+1 < len(s) && s[i+1] == '\n' {
				continue
			}
			lines = append(lines, string(current))
			current = current[:0]
		} else {
			current = append(current, s[i])
		}
	}

	// Add final line (even if empty after last newline)
	lines = append(lines, string(current))

	return lines
}

// joinLines joins lines with newline separators.
func joinLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	// Calculate total length
	total := 0
	for _, line := range lines {
		total += len(line)
	}
	total += len(lines) - 1 // newlines

	// Build result
	result := make([]byte, 0, total)
	for i, line := range lines {
		if i > 0 {
			result = append(result, '\n')
		}
		result = append(result, line...)
	}

	return string(result)
}

// isWordChar returns true if the byte is a word character.
func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}

// isWhitespace returns true if the byte is whitespace.
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t'
}
