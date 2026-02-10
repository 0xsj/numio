package rpc

// ════════════════════════════════════════════════════════════════
// INITIALIZATION
// ════════════════════════════════════════════════════════════════

// InitializeParams are the parameters for the initialize request.
type InitializeParams struct {
	// ClientName identifies the UI client (e.g., "numio-macos", "numio-windows")
	ClientName string `json:"clientName,omitempty"`

	// ClientVersion is the UI client version
	ClientVersion string `json:"clientVersion,omitempty"`

	// Viewport is the initial viewport size
	Viewport Viewport `json:"viewport"`
}

// InitializeResult is the response to an initialize request.
type InitializeResult struct {
	// CoreVersion is the numio-core version
	CoreVersion string `json:"coreVersion"`

	// ProtocolVersion is the RPC protocol version
	ProtocolVersion string `json:"protocolVersion"`

	// Capabilities lists supported features
	Capabilities Capabilities `json:"capabilities"`

	// InitialState is the initial render state
	InitialState RenderState `json:"initialState"`
}

// Capabilities describes what the core supports.
type Capabilities struct {
	// SupportsCurrency indicates live currency conversion is available
	SupportsCurrency bool `json:"supportsCurrency"`

	// SupportsCrypto indicates cryptocurrency prices are available
	SupportsCrypto bool `json:"supportsCrypto"`

	// SupportsMetals indicates precious metal prices are available
	SupportsMetals bool `json:"supportsMetals"`

	// SupportsUnits indicates unit conversion is available
	SupportsUnits bool `json:"supportsUnits"`

	// SupportsDates indicates date arithmetic is available
	SupportsDates bool `json:"supportsDates"`

	// SupportsExplain indicates explain mode is available
	SupportsExplain bool `json:"supportsExplain"`
}

// ════════════════════════════════════════════════════════════════
// VIEWPORT
// ════════════════════════════════════════════════════════════════

// Viewport represents the visible area dimensions.
type Viewport struct {
	// Width in characters
	Width int `json:"width"`

	// Height in lines
	Height int `json:"height"`
}

// ResizeParams are the parameters for a resize request.
type ResizeParams struct {
	Viewport Viewport `json:"viewport"`
}

// ════════════════════════════════════════════════════════════════
// INPUT EVENTS
// ════════════════════════════════════════════════════════════════

// KeyEvent represents a key press from the UI.
type KeyEvent struct {
	// Key is the key identifier
	// Single characters: "a", "A", "1", "!"
	// Special keys: "Enter", "Escape", "Backspace", "Delete", "Tab"
	// Arrow keys: "ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight"
	// Function keys: "F1", "F2", etc.
	Key string `json:"key"`

	// Modifiers are the modifier keys held during the press
	// Possible values: "ctrl", "alt", "shift", "cmd" (macOS), "meta" (Windows)
	Modifiers []string `json:"modifiers,omitempty"`
}

// HasModifier checks if a specific modifier was held.
func (k *KeyEvent) HasModifier(mod string) bool {
	for _, m := range k.Modifiers {
		if m == mod {
			return true
		}
	}
	return false
}

// IsCtrl checks if Ctrl was held.
func (k *KeyEvent) IsCtrl() bool {
	return k.HasModifier("ctrl")
}

// IsAlt checks if Alt/Option was held.
func (k *KeyEvent) IsAlt() bool {
	return k.HasModifier("alt")
}

// IsShift checks if Shift was held.
func (k *KeyEvent) IsShift() bool {
	return k.HasModifier("shift")
}

// IsCmd checks if Cmd (macOS) or Meta (Windows) was held.
func (k *KeyEvent) IsCmd() bool {
	return k.HasModifier("cmd") || k.HasModifier("meta")
}

// KeypressParams are the parameters for a keypress request.
type KeypressParams struct {
	Event KeyEvent `json:"event"`
}

// ════════════════════════════════════════════════════════════════
// RENDER STATE
// ════════════════════════════════════════════════════════════════

// RenderState represents the complete UI state to render.
type RenderState struct {
	// Lines are the editor lines with styled content
	Lines []RenderLine `json:"lines"`

	// Cursor is the current cursor position
	Cursor CursorState `json:"cursor"`

	// Mode is the current editor mode
	Mode EditorMode `json:"mode"`

	// StatusBar is the status bar content
	StatusBar *StatusBar `json:"statusBar,omitempty"`

	// Popup is an optional popup (help, explain, etc.)
	Popup *Popup `json:"popup,omitempty"`

	// Dirty indicates unsaved changes exist
	Dirty bool `json:"dirty"`
}

// RenderLine represents a single line in the editor.
type RenderLine struct {
	// Number is the 1-based line number
	Number int `json:"number"`

	// RelativeNumber is the distance from the cursor line (0 = cursor line)
	RelativeNumber int `json:"relativeNumber"`

	// Input contains the styled input spans
	Input []Span `json:"input"`

	// Result contains the styled result spans (right-aligned)
	Result []Span `json:"result,omitempty"`

	// IsCurrentLine indicates this is the cursor's line
	IsCurrentLine bool `json:"isCurrentLine,omitempty"`
}

// Span represents a styled text segment.
type Span struct {
	// Text is the content
	Text string `json:"text"`

	// Style is the style identifier
	Style SpanStyle `json:"style"`
}

// SpanStyle identifies the styling to apply.
type SpanStyle string

const (
	StyleDefault    SpanStyle = "default"
	StyleNumber     SpanStyle = "number"
	StyleOperator   SpanStyle = "operator"
	StylePercent    SpanStyle = "percent"
	StyleCurrency   SpanStyle = "currency"
	StyleUnit       SpanStyle = "unit"
	StyleMetal      SpanStyle = "metal"
	StyleCrypto     SpanStyle = "crypto"
	StyleFunction   SpanStyle = "function"
	StyleVariable   SpanStyle = "variable"
	StyleComment    SpanStyle = "comment"
	StyleString     SpanStyle = "string"
	StyleKeyword    SpanStyle = "keyword"
	StyleResult     SpanStyle = "result"
	StyleError      SpanStyle = "error"
	StylePending    SpanStyle = "pending"
	StyleLineNumber SpanStyle = "lineNumber"
	StyleSelection  SpanStyle = "selection"
)

// ════════════════════════════════════════════════════════════════
// CURSOR & SELECTION
// ════════════════════════════════════════════════════════════════

// CursorState represents the cursor position and visibility.
type CursorState struct {
	// Row is the 0-based row index
	Row int `json:"row"`

	// Col is the 0-based column index
	Col int `json:"col"`

	// Visible indicates if cursor should be shown
	Visible bool `json:"visible"`

	// Style is the cursor style based on mode
	Style CursorStyle `json:"style"`

	// Selection is the current selection range (if any)
	Selection *Selection `json:"selection,omitempty"`
}

// CursorStyle defines how the cursor should be rendered.
type CursorStyle string

const (
	CursorBlock     CursorStyle = "block"     // Normal mode
	CursorLine      CursorStyle = "line"      // Insert mode
	CursorUnderline CursorStyle = "underline" // Replace mode
)

// Selection represents a text selection range.
type Selection struct {
	// StartRow is the selection start row (0-based)
	StartRow int `json:"startRow"`

	// StartCol is the selection start column (0-based)
	StartCol int `json:"startCol"`

	// EndRow is the selection end row (0-based)
	EndRow int `json:"endRow"`

	// EndCol is the selection end column (0-based)
	EndCol int `json:"endCol"`
}

// ════════════════════════════════════════════════════════════════
// EDITOR MODE
// ════════════════════════════════════════════════════════════════

// EditorMode represents the current vim-like mode.
type EditorMode string

const (
	ModeNormal EditorMode = "normal"
	ModeInsert EditorMode = "insert"
	ModeVisual EditorMode = "visual"
)

// ════════════════════════════════════════════════════════════════
// STATUS BAR
// ════════════════════════════════════════════════════════════════

// StatusBar represents the status bar content.
type StatusBar struct {
	// Left is the left-aligned content
	Left []Span `json:"left,omitempty"`

	// Center is the center-aligned content
	Center []Span `json:"center,omitempty"`

	// Right is the right-aligned content
	Right []Span `json:"right,omitempty"`

	// Mode is displayed mode indicator
	Mode string `json:"mode,omitempty"`

	// Position shows cursor position (e.g., "3:14")
	Position string `json:"position,omitempty"`
}

// ════════════════════════════════════════════════════════════════
// POPUPS
// ════════════════════════════════════════════════════════════════

// Popup represents a modal popup overlay.
type Popup struct {
	// Type identifies the popup kind
	Type PopupType `json:"type"`

	// Title is the popup header
	Title string `json:"title,omitempty"`

	// Content is the popup body
	Content []PopupLine `json:"content"`
}

// PopupType identifies the kind of popup.
type PopupType string

const (
	PopupHelp    PopupType = "help"
	PopupExplain PopupType = "explain"
	PopupError   PopupType = "error"
)

// PopupLine is a line within a popup.
type PopupLine struct {
	Spans []Span `json:"spans"`
}

// ════════════════════════════════════════════════════════════════
// RATE STATUS
// ════════════════════════════════════════════════════════════════

// RateStatus represents the status of rate fetching.
type RateStatus struct {
	// State is the current fetch state
	State RateState `json:"state"`

	// Message is a human-readable status message
	Message string `json:"message,omitempty"`

	// UpdatedAt is when rates were last updated (Unix timestamp)
	UpdatedAt int64 `json:"updatedAt,omitempty"`

	// Error is the error message if state is "error"
	Error string `json:"error,omitempty"`
}

// RateState represents the rate fetching state.
type RateState string

const (
	RateStateIdle     RateState = "idle"
	RateStateFetching RateState = "fetching"
	RateStateSuccess  RateState = "success"
	RateStateError    RateState = "error"
)

// RateStatusParams are the parameters for rate status notifications.
type RateStatusParams struct {
	Status RateStatus `json:"status"`
}

// ════════════════════════════════════════════════════════════════
// OPTIONS
// ════════════════════════════════════════════════════════════════

// SetOptionParams are the parameters for setOption requests.
type SetOptionParams struct {
	// Option is the option name
	Option string `json:"option"`

	// Value is the option value (type depends on option)
	Value any `json:"value"`
}

// Known option names.
const (
	OptionPrecision = "precision" // int: decimal places
	OptionTheme     = "theme"     // string: theme name
	OptionVimMode   = "vimMode"   // bool: enable vim bindings
)

// ════════════════════════════════════════════════════════════════
// EXPLAIN
// ════════════════════════════════════════════════════════════════

// ExplainResult contains step-by-step explanation of a calculation.
type ExplainResult struct {
	// Input is the original expression
	Input string `json:"input"`

	// Steps are the calculation steps
	Steps []ExplainStep `json:"steps"`

	// FinalResult is the final computed value
	FinalResult []Span `json:"finalResult"`
}

// ExplainStep is a single step in the explanation.
type ExplainStep struct {
	// Description explains what this step does
	Description string `json:"description"`

	// Expression shows the expression at this step
	Expression []Span `json:"expression"`

	// Result shows the step result
	Result []Span `json:"result,omitempty"`
}
