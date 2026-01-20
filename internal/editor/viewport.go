package editor

// ════════════════════════════════════════════════════════════════
// VIEWPORT
// ════════════════════════════════════════════════════════════════

// Viewport represents the visible region of the editor.
type Viewport struct {
	// Top-left position of the viewport in the buffer.
	topRow  int
	leftCol int

	// Visible dimensions.
	width  int
	height int

	// Scroll margins (lines/cols to keep visible around cursor).
	scrollOffVertical   int
	scrollOffHorizontal int
}

// NewViewport creates a viewport with the given dimensions.
func NewViewport(width, height int) *Viewport {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	return &Viewport{
		topRow:              0,
		leftCol:             0,
		width:               width,
		height:              height,
		scrollOffVertical:   3,
		scrollOffHorizontal: 5,
	}
}

// ════════════════════════════════════════════════════════════════
// ACCESSORS
// ════════════════════════════════════════════════════════════════

// TopRow returns the first visible row.
func (v *Viewport) TopRow() int {
	return v.topRow
}

// LeftCol returns the first visible column.
func (v *Viewport) LeftCol() int {
	return v.leftCol
}

// Width returns the viewport width in columns.
func (v *Viewport) Width() int {
	return v.width
}

// Height returns the viewport height in rows.
func (v *Viewport) Height() int {
	return v.height
}

// BottomRow returns the last visible row (exclusive).
func (v *Viewport) BottomRow() int {
	return v.topRow + v.height
}

// RightCol returns the last visible column (exclusive).
func (v *Viewport) RightCol() int {
	return v.leftCol + v.width
}

// VisibleRange returns the visible row range (startRow, endRow exclusive).
func (v *Viewport) VisibleRange() (int, int) {
	return v.topRow, v.topRow + v.height
}

// ════════════════════════════════════════════════════════════════
// SETTERS
// ════════════════════════════════════════════════════════════════

// SetSize updates the viewport dimensions.
func (v *Viewport) SetSize(width, height int) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	v.width = width
	v.height = height
}

// SetTopRow sets the first visible row.
func (v *Viewport) SetTopRow(row int) {
	if row < 0 {
		row = 0
	}
	v.topRow = row
}

// SetLeftCol sets the first visible column.
func (v *Viewport) SetLeftCol(col int) {
	if col < 0 {
		col = 0
	}
	v.leftCol = col
}

// SetScrollOff sets the scroll margins.
func (v *Viewport) SetScrollOff(vertical, horizontal int) {
	if vertical < 0 {
		vertical = 0
	}
	if horizontal < 0 {
		horizontal = 0
	}
	v.scrollOffVertical = vertical
	v.scrollOffHorizontal = horizontal
}

// ScrollOffVertical returns the vertical scroll margin.
func (v *Viewport) ScrollOffVertical() int {
	return v.scrollOffVertical
}

// ScrollOffHorizontal returns the horizontal scroll margin.
func (v *Viewport) ScrollOffHorizontal() int {
	return v.scrollOffHorizontal
}

// ════════════════════════════════════════════════════════════════
// SCROLLING
// ════════════════════════════════════════════════════════════════

// ScrollTo scrolls to make the given position visible.
// Returns true if the viewport moved.
func (v *Viewport) ScrollTo(row, col int) bool {
	moved := false

	// Vertical scrolling
	if row < v.topRow+v.scrollOffVertical {
		v.topRow = row - v.scrollOffVertical
		moved = true
	} else if row >= v.topRow+v.height-v.scrollOffVertical {
		v.topRow = row - v.height + v.scrollOffVertical + 1
		moved = true
	}

	// Horizontal scrolling
	if col < v.leftCol+v.scrollOffHorizontal {
		v.leftCol = col - v.scrollOffHorizontal
		moved = true
	} else if col >= v.leftCol+v.width-v.scrollOffHorizontal {
		v.leftCol = col - v.width + v.scrollOffHorizontal + 1
		moved = true
	}

	// Clamp to valid ranges
	if v.topRow < 0 {
		v.topRow = 0
	}
	if v.leftCol < 0 {
		v.leftCol = 0
	}

	return moved
}

// ScrollToCenter centers the viewport on the given position.
func (v *Viewport) ScrollToCenter(row, col int) {
	v.topRow = row - v.height/2
	v.leftCol = col - v.width/2

	if v.topRow < 0 {
		v.topRow = 0
	}
	if v.leftCol < 0 {
		v.leftCol = 0
	}
}

// ScrollUp scrolls the viewport up by n rows.
func (v *Viewport) ScrollUp(n int) {
	v.topRow -= n
	if v.topRow < 0 {
		v.topRow = 0
	}
}

// ScrollDown scrolls the viewport down by n rows.
// maxRow is the maximum valid top row (total lines - height).
func (v *Viewport) ScrollDown(n int, maxRow int) {
	v.topRow += n
	if maxRow < 0 {
		maxRow = 0
	}
	if v.topRow > maxRow {
		v.topRow = maxRow
	}
}

// ScrollLeft scrolls the viewport left by n columns.
func (v *Viewport) ScrollLeft(n int) {
	v.leftCol -= n
	if v.leftCol < 0 {
		v.leftCol = 0
	}
}

// ScrollRight scrolls the viewport right by n columns.
func (v *Viewport) ScrollRight(n int) {
	v.leftCol += n
}

// PageUp scrolls up by one page (viewport height).
func (v *Viewport) PageUp() {
	v.ScrollUp(v.height)
}

// PageDown scrolls down by one page.
func (v *Viewport) PageDown(maxRow int) {
	v.ScrollDown(v.height, maxRow)
}

// HalfPageUp scrolls up by half a page.
func (v *Viewport) HalfPageUp() {
	v.ScrollUp(v.height / 2)
}

// HalfPageDown scrolls down by half a page.
func (v *Viewport) HalfPageDown(maxRow int) {
	v.ScrollDown(v.height/2, maxRow)
}

// ════════════════════════════════════════════════════════════════
// QUERIES
// ════════════════════════════════════════════════════════════════

// IsVisible returns true if the given position is visible.
func (v *Viewport) IsVisible(row, col int) bool {
	return row >= v.topRow &&
		row < v.topRow+v.height &&
		col >= v.leftCol &&
		col < v.leftCol+v.width
}

// IsRowVisible returns true if the given row is visible.
func (v *Viewport) IsRowVisible(row int) bool {
	return row >= v.topRow && row < v.topRow+v.height
}

// IsColVisible returns true if the given column is visible.
func (v *Viewport) IsColVisible(col int) bool {
	return col >= v.leftCol && col < v.leftCol+v.width
}

// ToViewCoords converts buffer coordinates to viewport coordinates.
// Returns (-1, -1) if not visible.
func (v *Viewport) ToViewCoords(row, col int) (int, int) {
	if !v.IsVisible(row, col) {
		return -1, -1
	}
	return row - v.topRow, col - v.leftCol
}

// ToBufferCoords converts viewport coordinates to buffer coordinates.
func (v *Viewport) ToBufferCoords(viewRow, viewCol int) (int, int) {
	return viewRow + v.topRow, viewCol + v.leftCol
}

// ════════════════════════════════════════════════════════════════
// CLAMPING
// ════════════════════════════════════════════════════════════════

// Clamp ensures the viewport is within valid bounds.
// totalLines is the total number of lines in the buffer.
func (v *Viewport) Clamp(totalLines int) {
	if v.topRow < 0 {
		v.topRow = 0
	}
	if v.leftCol < 0 {
		v.leftCol = 0
	}

	// Don't scroll past the last line
	maxTopRow := totalLines - v.height
	if maxTopRow < 0 {
		maxTopRow = 0
	}
	if v.topRow > maxTopRow {
		v.topRow = maxTopRow
	}
}

// ClampToBuffer ensures the viewport doesn't extend past the buffer.
func (v *Viewport) ClampToBuffer(totalLines, maxLineLength int) {
	v.Clamp(totalLines)

	maxLeftCol := maxLineLength - v.width
	if maxLeftCol < 0 {
		maxLeftCol = 0
	}
	if v.leftCol > maxLeftCol {
		v.leftCol = maxLeftCol
	}
}

// ════════════════════════════════════════════════════════════════
// VISIBLE LINES
// ════════════════════════════════════════════════════════════════

// VisibleLines returns the range of visible line indices.
// Clamped to the actual buffer size.
func (v *Viewport) VisibleLines(totalLines int) (start, end int) {
	start = v.topRow
	end = v.topRow + v.height

	if start < 0 {
		start = 0
	}
	if end > totalLines {
		end = totalLines
	}

	return start, end
}

// VisibleLineCount returns the number of visible lines.
func (v *Viewport) VisibleLineCount(totalLines int) int {
	start, end := v.VisibleLines(totalLines)
	return end - start
}

// ════════════════════════════════════════════════════════════════
// LINE POSITIONS
// ════════════════════════════════════════════════════════════════

// CursorViewRow returns the cursor's row in viewport coordinates.
// Returns -1 if not visible.
func (v *Viewport) CursorViewRow(cursorRow int) int {
	if cursorRow < v.topRow || cursorRow >= v.topRow+v.height {
		return -1
	}
	return cursorRow - v.topRow
}

// CursorViewCol returns the cursor's column in viewport coordinates.
// Returns -1 if not visible.
func (v *Viewport) CursorViewCol(cursorCol int) int {
	if cursorCol < v.leftCol || cursorCol >= v.leftCol+v.width {
		return -1
	}
	return cursorCol - v.leftCol
}

// ════════════════════════════════════════════════════════════════
// ENSURE VISIBLE
// ════════════════════════════════════════════════════════════════

// EnsureCursorVisible scrolls minimally to make the cursor visible.
// Returns true if the viewport moved.
func (v *Viewport) EnsureCursorVisible(cursorRow, cursorCol int) bool {
	return v.ScrollTo(cursorRow, cursorCol)
}

// EnsureRowVisible scrolls minimally to make a row visible.
func (v *Viewport) EnsureRowVisible(row int) bool {
	moved := false

	if row < v.topRow {
		v.topRow = row
		moved = true
	} else if row >= v.topRow+v.height {
		v.topRow = row - v.height + 1
		moved = true
	}

	if v.topRow < 0 {
		v.topRow = 0
	}

	return moved
}

// ════════════════════════════════════════════════════════════════
// CLONING
// ════════════════════════════════════════════════════════════════

// Clone creates a copy of the viewport.
func (v *Viewport) Clone() *Viewport {
	return &Viewport{
		topRow:              v.topRow,
		leftCol:             v.leftCol,
		width:               v.width,
		height:              v.height,
		scrollOffVertical:   v.scrollOffVertical,
		scrollOffHorizontal: v.scrollOffHorizontal,
	}
}
