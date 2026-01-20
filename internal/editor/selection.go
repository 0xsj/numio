package editor

// ════════════════════════════════════════════════════════════════
// SELECTION
// ════════════════════════════════════════════════════════════════

// Selection represents a text selection in the editor.
type Selection struct {
	// Anchor is where the selection started.
	anchorRow int
	anchorCol int

	// Active is the current cursor position (movable end).
	activeRow int
	activeCol int

	// Mode determines selection behavior.
	mode SelectionMode
}

// SelectionMode determines how the selection behaves.
type SelectionMode int

const (
	// SelectionChar is character-wise selection (v in vim).
	SelectionChar SelectionMode = iota

	// SelectionLine is line-wise selection (V in vim).
	SelectionLine

	// SelectionBlock is block/column selection (Ctrl+v in vim).
	SelectionBlock
)

// NewSelection creates a character selection starting at the given position.
func NewSelection(row, col int) *Selection {
	return &Selection{
		anchorRow: row,
		anchorCol: col,
		activeRow: row,
		activeCol: col,
		mode:      SelectionChar,
	}
}

// NewSelectionWithMode creates a selection with the specified mode.
func NewSelectionWithMode(row, col int, mode SelectionMode) *Selection {
	return &Selection{
		anchorRow: row,
		anchorCol: col,
		activeRow: row,
		activeCol: col,
		mode:      mode,
	}
}

// ════════════════════════════════════════════════════════════════
// ACCESSORS
// ════════════════════════════════════════════════════════════════

// Anchor returns the anchor position.
func (s *Selection) Anchor() (int, int) {
	return s.anchorRow, s.anchorCol
}

// Active returns the active (cursor) position.
func (s *Selection) Active() (int, int) {
	return s.activeRow, s.activeCol
}

// Mode returns the selection mode.
func (s *Selection) Mode() SelectionMode {
	return s.mode
}

// Start returns the start position (minimum of anchor and active).
func (s *Selection) Start() (int, int) {
	if s.anchorRow < s.activeRow {
		return s.anchorRow, s.anchorCol
	}
	if s.anchorRow > s.activeRow {
		return s.activeRow, s.activeCol
	}
	// Same row
	if s.anchorCol < s.activeCol {
		return s.anchorRow, s.anchorCol
	}
	return s.activeRow, s.activeCol
}

// End returns the end position (maximum of anchor and active).
func (s *Selection) End() (int, int) {
	if s.anchorRow > s.activeRow {
		return s.anchorRow, s.anchorCol
	}
	if s.anchorRow < s.activeRow {
		return s.activeRow, s.activeCol
	}
	// Same row
	if s.anchorCol > s.activeCol {
		return s.anchorRow, s.anchorCol
	}
	return s.activeRow, s.activeCol
}

// Bounds returns the normalized selection bounds.
// Returns (startRow, startCol, endRow, endCol).
func (s *Selection) Bounds() (int, int, int, int) {
	startRow, startCol := s.Start()
	endRow, endCol := s.End()
	return startRow, startCol, endRow, endCol
}

// StartRow returns the starting row of the selection.
func (s *Selection) StartRow() int {
	if s.anchorRow < s.activeRow {
		return s.anchorRow
	}
	return s.activeRow
}

// EndRow returns the ending row of the selection.
func (s *Selection) EndRow() int {
	if s.anchorRow > s.activeRow {
		return s.anchorRow
	}
	return s.activeRow
}

// ════════════════════════════════════════════════════════════════
// SETTERS
// ════════════════════════════════════════════════════════════════

// SetActive updates the active (cursor) position.
func (s *Selection) SetActive(row, col int) {
	s.activeRow = row
	s.activeCol = col
}

// SetAnchor updates the anchor position.
func (s *Selection) SetAnchor(row, col int) {
	s.anchorRow = row
	s.anchorCol = col
}

// SetMode changes the selection mode.
func (s *Selection) SetMode(mode SelectionMode) {
	s.mode = mode
}

// ════════════════════════════════════════════════════════════════
// QUERIES
// ════════════════════════════════════════════════════════════════

// IsEmpty returns true if the selection has no extent.
func (s *Selection) IsEmpty() bool {
	return s.anchorRow == s.activeRow && s.anchorCol == s.activeCol
}

// IsForward returns true if active is after anchor.
func (s *Selection) IsForward() bool {
	if s.activeRow > s.anchorRow {
		return true
	}
	if s.activeRow < s.anchorRow {
		return false
	}
	return s.activeCol >= s.anchorCol
}

// IsReversed returns true if active is before anchor.
func (s *Selection) IsReversed() bool {
	return !s.IsForward()
}

// Contains returns true if the given position is within the selection.
func (s *Selection) Contains(row, col int) bool {
	startRow, startCol, endRow, endCol := s.Bounds()

	switch s.mode {
	case SelectionChar:
		return s.containsChar(row, col, startRow, startCol, endRow, endCol)

	case SelectionLine:
		return row >= startRow && row <= endRow

	case SelectionBlock:
		return s.containsBlock(row, col, startRow, startCol, endRow, endCol)

	default:
		return false
	}
}

func (s *Selection) containsChar(row, col, startRow, startCol, endRow, endCol int) bool {
	if row < startRow || row > endRow {
		return false
	}

	if row == startRow && row == endRow {
		// Single line selection
		return col >= startCol && col < endCol
	}

	if row == startRow {
		return col >= startCol
	}

	if row == endRow {
		return col < endCol
	}

	// Middle row - entire row is selected
	return true
}

func (s *Selection) containsBlock(row, col, startRow, startCol, endRow, endCol int) bool {
	if row < startRow || row > endRow {
		return false
	}

	minCol := startCol
	maxCol := endCol
	if minCol > maxCol {
		minCol, maxCol = maxCol, minCol
	}

	return col >= minCol && col < maxCol
}

// ContainsRow returns true if the given row is part of the selection.
func (s *Selection) ContainsRow(row int) bool {
	startRow, _, endRow, _ := s.Bounds()
	return row >= startRow && row <= endRow
}

// ════════════════════════════════════════════════════════════════
// LINE-WISE HELPERS
// ════════════════════════════════════════════════════════════════

// LineBounds returns the selection bounds expanded to full lines.
// Useful for line-wise operations.
func (s *Selection) LineBounds() (int, int) {
	startRow := s.StartRow()
	endRow := s.EndRow()
	return startRow, endRow
}

// ExpandToLines converts the selection to line-wise mode.
func (s *Selection) ExpandToLines() {
	s.mode = SelectionLine
	s.anchorCol = 0
	s.activeCol = 0
}

// ════════════════════════════════════════════════════════════════
// BLOCK SELECTION HELPERS
// ════════════════════════════════════════════════════════════════

// BlockBounds returns the column bounds for block selection.
// Returns (minCol, maxCol).
func (s *Selection) BlockBounds() (int, int) {
	minCol := s.anchorCol
	maxCol := s.activeCol
	if minCol > maxCol {
		minCol, maxCol = maxCol, minCol
	}
	return minCol, maxCol
}

// BlockContainsCol returns true if the column is within the block selection.
func (s *Selection) BlockContainsCol(col int) bool {
	minCol, maxCol := s.BlockBounds()
	return col >= minCol && col < maxCol
}

// ════════════════════════════════════════════════════════════════
// MANIPULATION
// ════════════════════════════════════════════════════════════════

// Swap swaps the anchor and active positions.
func (s *Selection) Swap() {
	s.anchorRow, s.activeRow = s.activeRow, s.anchorRow
	s.anchorCol, s.activeCol = s.activeCol, s.anchorCol
}

// Normalize ensures anchor is before active.
func (s *Selection) Normalize() {
	if s.IsReversed() {
		s.Swap()
	}
}

// Extend extends the selection to include the given position.
func (s *Selection) Extend(row, col int) {
	s.activeRow = row
	s.activeCol = col
}

// ExtendByRows extends the selection by n rows.
func (s *Selection) ExtendByRows(n int) {
	s.activeRow += n
	if s.activeRow < 0 {
		s.activeRow = 0
	}
}

// ExtendByCols extends the selection by n columns.
func (s *Selection) ExtendByCols(n int) {
	s.activeCol += n
	if s.activeCol < 0 {
		s.activeCol = 0
	}
}

// ════════════════════════════════════════════════════════════════
// CLONING
// ════════════════════════════════════════════════════════════════

// Clone creates a copy of the selection.
func (s *Selection) Clone() *Selection {
	return &Selection{
		anchorRow: s.anchorRow,
		anchorCol: s.anchorCol,
		activeRow: s.activeRow,
		activeCol: s.activeCol,
		mode:      s.mode,
	}
}

// ════════════════════════════════════════════════════════════════
// SELECTION RANGE (for iteration)
// ════════════════════════════════════════════════════════════════

// Range represents a range on a single line.
type Range struct {
	Row      int
	StartCol int
	EndCol   int // Exclusive
}

// Ranges returns the selection as a slice of line ranges.
// Useful for rendering and text extraction.
func (s *Selection) Ranges(lineLengths []int) []Range {
	startRow, startCol, endRow, endCol := s.Bounds()

	// Clamp to available lines
	if startRow >= len(lineLengths) {
		return nil
	}
	if endRow >= len(lineLengths) {
		endRow = len(lineLengths) - 1
		endCol = lineLengths[endRow]
	}

	switch s.mode {
	case SelectionChar:
		return s.charRanges(startRow, startCol, endRow, endCol, lineLengths)

	case SelectionLine:
		return s.lineRanges(startRow, endRow, lineLengths)

	case SelectionBlock:
		return s.blockRanges(startRow, endRow, lineLengths)

	default:
		return nil
	}
}

func (s *Selection) charRanges(startRow, startCol, endRow, endCol int, lineLengths []int) []Range {
	var ranges []Range

	for row := startRow; row <= endRow; row++ {
		lineLen := lineLengths[row]
		r := Range{Row: row}

		if row == startRow && row == endRow {
			// Single line
			r.StartCol = startCol
			r.EndCol = endCol
		} else if row == startRow {
			// First line
			r.StartCol = startCol
			r.EndCol = lineLen
		} else if row == endRow {
			// Last line
			r.StartCol = 0
			r.EndCol = endCol
		} else {
			// Middle line
			r.StartCol = 0
			r.EndCol = lineLen
		}

		// Clamp
		if r.StartCol > lineLen {
			r.StartCol = lineLen
		}
		if r.EndCol > lineLen {
			r.EndCol = lineLen
		}

		ranges = append(ranges, r)
	}

	return ranges
}

func (s *Selection) lineRanges(startRow, endRow int, lineLengths []int) []Range {
	var ranges []Range

	for row := startRow; row <= endRow; row++ {
		lineLen := lineLengths[row]
		ranges = append(ranges, Range{
			Row:      row,
			StartCol: 0,
			EndCol:   lineLen,
		})
	}

	return ranges
}

func (s *Selection) blockRanges(startRow, endRow int, lineLengths []int) []Range {
	minCol, maxCol := s.BlockBounds()
	var ranges []Range

	for row := startRow; row <= endRow; row++ {
		lineLen := lineLengths[row]

		r := Range{
			Row:      row,
			StartCol: minCol,
			EndCol:   maxCol,
		}

		// Clamp to line length
		if r.StartCol > lineLen {
			r.StartCol = lineLen
		}
		if r.EndCol > lineLen {
			r.EndCol = lineLen
		}

		ranges = append(ranges, r)
	}

	return ranges
}
