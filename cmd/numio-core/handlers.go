package main

import (
	"context"

	"github.com/0xsj/numio/internal/editor"
	"github.com/0xsj/numio/internal/rpc"
)

// ════════════════════════════════════════════════════════════════
// CORE HANDLER
// ════════════════════════════════════════════════════════════════

// CoreHandler handles RPC methods for the numio desktop backend.
type CoreHandler struct {
	editor  *editor.Editor
	server  *rpc.Server
	version string

	initialized bool
}

// NewCoreHandler creates a new core handler.
func NewCoreHandler(version string) *CoreHandler {
	return &CoreHandler{
		editor:      editor.NewEditor(),
		version:     version,
		initialized: false,
	}
}

// SetServer sets the RPC server reference (for sending notifications).
func (h *CoreHandler) SetServer(server *rpc.Server) {
	h.server = server
}

// Editor returns the editor instance.
func (h *CoreHandler) Editor() *editor.Editor {
	return h.editor
}

// RegisterHandlers registers all RPC method handlers.
func (h *CoreHandler) RegisterHandlers(registry *rpc.Registry) {
	// Request handlers
	registry.Register(rpc.MethodInitialize, rpc.TypedHandler(h.handleInitialize))
	registry.Register(rpc.MethodShutdown, rpc.TypedHandlerNoParams(h.handleShutdown))
	registry.Register(rpc.MethodKeypress, rpc.TypedHandler(h.handleKeypress))
	registry.Register(rpc.MethodResize, rpc.TypedHandler(h.handleResize))
	registry.Register(rpc.MethodGetState, rpc.TypedHandlerNoParams(h.handleGetState))
	registry.Register(rpc.MethodSetOption, rpc.TypedHandler(h.handleSetOption))
}

// ════════════════════════════════════════════════════════════════
// INITIALIZE
// ════════════════════════════════════════════════════════════════

func (h *CoreHandler) handleInitialize(ctx context.Context, params rpc.InitializeParams) (*rpc.InitializeResult, *rpc.Error) {
	// Set viewport size
	h.editor.Resize(params.Viewport.Width, params.Viewport.Height)

	h.initialized = true

	// Return initial state
	return &rpc.InitializeResult{
		CoreVersion:     h.version,
		ProtocolVersion: rpc.JSONRPCVersion,
		Capabilities: rpc.Capabilities{
			SupportsCurrency: true,
			SupportsCrypto:   true,
			SupportsMetals:   true,
			SupportsUnits:    true,
			SupportsDates:    true,
			SupportsExplain:  true,
		},
		InitialState: *h.editor.Render(),
	}, nil
}

// ════════════════════════════════════════════════════════════════
// SHUTDOWN
// ════════════════════════════════════════════════════════════════

func (h *CoreHandler) handleShutdown(ctx context.Context) (*struct{}, *rpc.Error) {
	// Perform cleanup if needed
	h.initialized = false

	// Signal server to stop (will be handled by main)
	if h.server != nil {
		go h.server.Stop()
	}

	return &struct{}{}, nil
}

// ════════════════════════════════════════════════════════════════
// KEYPRESS
// ════════════════════════════════════════════════════════════════

func (h *CoreHandler) handleKeypress(ctx context.Context, params rpc.KeypressParams) (*rpc.RenderState, *rpc.Error) {
	if !h.initialized {
		return nil, rpc.NotReadyError("not initialized")
	}

	// Process the key event
	h.processKey(params.Event)

	// Return updated render state
	state := h.editor.Render()
	return state, nil
}

func (h *CoreHandler) processKey(event rpc.KeyEvent) {
	mode := h.editor.Mode()

	// Handle mode-specific keys
	switch mode {
	case editor.ModeNormal:
		h.processNormalModeKey(event)
	case editor.ModeInsert:
		h.processInsertModeKey(event)
	case editor.ModeVisual:
		h.processVisualModeKey(event)
	}
}

func (h *CoreHandler) processNormalModeKey(event rpc.KeyEvent) {
	key := event.Key

	// Handle modified keys first
	if event.IsCtrl() {
		switch key {
		case "r":
			h.editor.Redo()
		case "u":
			// Half page up
			h.editor.MoveUp(h.editor.Viewport().Height() / 2)
		case "d":
			// Half page down
			h.editor.MoveDown(h.editor.Viewport().Height() / 2)
		}
		return
	}

	// Normal mode keys
	switch key {
	// Mode switching
	case "i":
		h.editor.EnterInsertMode()
	case "I":
		h.editor.EnterInsertModeLineStart()
	case "a":
		h.editor.EnterInsertModeAppend()
	case "A":
		h.editor.EnterInsertModeLineEnd()
	case "o":
		h.editor.OpenLineBelow()
	case "O":
		h.editor.OpenLineAbove()
	case "v":
		h.editor.EnterVisualMode()
	case "V":
		h.editor.EnterVisualLineMode()

	// Movement
	case "h", "ArrowLeft":
		h.editor.MoveLeft(1)
	case "j", "ArrowDown":
		h.editor.MoveDown(1)
	case "k", "ArrowUp":
		h.editor.MoveUp(1)
	case "l", "ArrowRight":
		h.editor.MoveRight(1)
	case "w":
		h.editor.MoveWordForward(1)
	case "b":
		h.editor.MoveWordBackward(1)
	case "0":
		h.editor.MoveToLineStart()
	case "$":
		h.editor.MoveToLineEnd()
	case "^":
		h.editor.MoveToFirstNonBlank()
	case "g":
		// Would need to track for 'gg', simplified here
		h.editor.MoveToTop()
	case "G":
		h.editor.MoveToBottom()

	// Editing
	case "x":
		h.editor.DeleteChar()
	case "X":
		h.editor.DeleteCharBack()
	case "d":
		// Would need to track for 'dd', 'dw', etc.
		// Simplified: delete line
		h.editor.DeleteLine()
	case "D":
		h.editor.DeleteToLineEnd()
	case "y":
		// Simplified: yank line
		h.editor.YankLine()
	case "p":
		h.editor.Paste()
	case "P":
		h.editor.PasteBefore()
	case "u":
		h.editor.Undo()
	case "J":
		h.editor.JoinLines()
	}
}

func (h *CoreHandler) processInsertModeKey(event rpc.KeyEvent) {
	key := event.Key

	// Handle escape to exit insert mode
	if key == "Escape" {
		h.editor.EnterNormalMode()
		return
	}

	// Handle control keys
	if event.IsCtrl() {
		switch key {
		case "c":
			h.editor.EnterNormalMode()
		case "w":
			h.editor.DeleteWord()
		case "u":
			h.editor.DeleteToLineEnd() // Actually should delete to line start
		}
		return
	}

	// Handle special keys
	switch key {
	case "Backspace":
		h.editor.DeleteCharBack()
	case "Delete":
		h.editor.DeleteChar()
	case "Enter":
		h.editor.InsertNewline()
	case "Tab":
		h.editor.InsertString("  ")
	case "ArrowLeft":
		h.editor.MoveLeft(1)
	case "ArrowRight":
		h.editor.MoveRight(1)
	case "ArrowUp":
		h.editor.MoveUp(1)
	case "ArrowDown":
		h.editor.MoveDown(1)
	default:
		// Insert regular character
		if len(key) == 1 {
			h.editor.InsertChar(rune(key[0]))
		}
	}
}

func (h *CoreHandler) processVisualModeKey(event rpc.KeyEvent) {
	key := event.Key

	// Escape exits visual mode
	if key == "Escape" {
		h.editor.EnterNormalMode()
		return
	}

	// Movement keys extend selection
	switch key {
	case "h", "ArrowLeft":
		h.editor.MoveLeft(1)
	case "j", "ArrowDown":
		h.editor.MoveDown(1)
	case "k", "ArrowUp":
		h.editor.MoveUp(1)
	case "l", "ArrowRight":
		h.editor.MoveRight(1)
	case "w":
		h.editor.MoveWordForward(1)
	case "b":
		h.editor.MoveWordBackward(1)
	case "0":
		h.editor.MoveToLineStart()
	case "$":
		h.editor.MoveToLineEnd()
	case "G":
		h.editor.MoveToBottom()
	case "g":
		h.editor.MoveToTop()

	// Actions on selection
	case "y":
		h.editor.YankSelection()
		h.editor.EnterNormalMode()
	case "d", "x":
		h.editor.YankSelection()
		// TODO: delete selection
		h.editor.EnterNormalMode()
	case "v":
		// Toggle back to normal
		h.editor.EnterNormalMode()
	}
}

// ════════════════════════════════════════════════════════════════
// RESIZE
// ════════════════════════════════════════════════════════════════

func (h *CoreHandler) handleResize(ctx context.Context, params rpc.ResizeParams) (*rpc.RenderState, *rpc.Error) {
	if !h.initialized {
		return nil, rpc.NotReadyError("not initialized")
	}

	h.editor.Resize(params.Viewport.Width, params.Viewport.Height)

	state := h.editor.Render()
	return state, nil
}

// ════════════════════════════════════════════════════════════════
// GET STATE
// ════════════════════════════════════════════════════════════════

func (h *CoreHandler) handleGetState(ctx context.Context) (*rpc.RenderState, *rpc.Error) {
	if !h.initialized {
		return nil, rpc.NotReadyError("not initialized")
	}

	state := h.editor.Render()
	return state, nil
}

// ════════════════════════════════════════════════════════════════
// SET OPTION
// ════════════════════════════════════════════════════════════════

func (h *CoreHandler) handleSetOption(ctx context.Context, params rpc.SetOptionParams) (*struct{}, *rpc.Error) {
	if !h.initialized {
		return nil, rpc.NotReadyError("not initialized")
	}

	switch params.Option {
	case rpc.OptionPrecision:
		if precision, ok := params.Value.(float64); ok {
			h.editor.Engine().SetPrecision(int(precision))
		} else {
			return nil, rpc.InvalidParamsError("precision must be a number")
		}

	case rpc.OptionTheme:
		// TODO: implement theme switching
		return nil, rpc.InvalidParamsError("theme not yet supported")

	case rpc.OptionVimMode:
		// Always vim mode for now
		return nil, rpc.InvalidParamsError("vim mode cannot be disabled")

	default:
		return nil, rpc.InvalidParamsError("unknown option: " + params.Option)
	}

	return &struct{}{}, nil
}
