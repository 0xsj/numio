package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
)

// ════════════════════════════════════════════════════════════════
// APP
// ════════════════════════════════════════════════════════════════

// App is the main numio desktop application.
type App struct {
	fyneApp    fyne.App
	window     fyne.Window
	editor     *EditorWidget
	cursorTick *time.Ticker
	done       chan struct{}
}

// NewApp creates a new numio desktop application.
func NewApp() *App {
	return &App{
		done: make(chan struct{}),
	}
}

// Run starts the application.
func (a *App) Run() {
	// Create Fyne app
	a.fyneApp = app.New()
	a.fyneApp.Settings().SetTheme(&NumioTheme{})

	// Create main window
	a.window = a.fyneApp.NewWindow("Numio")
	a.window.Resize(fyne.NewSize(700, 500))
	a.window.CenterOnScreen()

	// Create editor widget
	a.editor = NewEditorWidget()

	// Set up layout - no padding for clean look
	content := container.NewWithoutLayout(a.editor)
	a.editor.Resize(a.window.Canvas().Size())

	a.window.SetContent(content)

	// Register keyboard shortcuts (Cmd/Ctrl + key)
	a.registerShortcuts()

	// Resize editor when window resizes
	a.window.Canvas().SetContent(content)
	go func() {
		// Give window time to initialize
		time.Sleep(50 * time.Millisecond)
		fyne.Do(func() {
			size := a.window.Canvas().Size()
			a.editor.Resize(size)
			a.editor.Refresh()
		})
	}()

	// Listen for resize events
	a.window.SetOnClosed(func() {
		a.Stop()
	})

	// Focus the editor
	a.window.Canvas().Focus(a.editor)

	// Start cursor blink
	a.startCursorBlink()

	// Set up window resize handler
	a.setupResizeHandler()

	// Show and run
	a.window.ShowAndRun()
}

// Stop stops the application.
func (a *App) Stop() {
	close(a.done)
	if a.cursorTick != nil {
		a.cursorTick.Stop()
	}
}

// ════════════════════════════════════════════════════════════════
// KEYBOARD SHORTCUTS
// ════════════════════════════════════════════════════════════════

func (a *App) registerShortcuts() {
	c := a.window.Canvas()

	// Cmd+/ and Ctrl+/ → Toggle help
	c.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeySlash,
		Modifier: fyne.KeyModifierSuper,
	}, func(_ fyne.Shortcut) {
		a.editor.ToggleHelp()
	})
	c.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeySlash,
		Modifier: fyne.KeyModifierControl,
	}, func(_ fyne.Shortcut) {
		a.editor.ToggleHelp()
	})

	// Cmd+E and Ctrl+E → Explain
	c.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyE,
		Modifier: fyne.KeyModifierSuper,
	}, func(_ fyne.Shortcut) {
		a.editor.ShowExplain()
	})
	c.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyE,
		Modifier: fyne.KeyModifierControl,
	}, func(_ fyne.Shortcut) {
		a.editor.ShowExplain()
	})

	// Cmd+K and Ctrl+K → Toggle vim mode
	c.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyK,
		Modifier: fyne.KeyModifierSuper,
	}, func(_ fyne.Shortcut) {
		a.editor.ToggleVimMode()
	})
	c.AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyK,
		Modifier: fyne.KeyModifierControl,
	}, func(_ fyne.Shortcut) {
		a.editor.ToggleVimMode()
	})
}

// ════════════════════════════════════════════════════════════════
// RESIZE HANDLER
// ════════════════════════════════════════════════════════════════

func (a *App) setupResizeHandler() {
	// Poll for size changes
	go func() {
		var lastSize fyne.Size
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fyne.Do(func() {
					size := a.window.Canvas().Size()
					if size != lastSize && size.Width > 0 && size.Height > 0 {
						lastSize = size
						a.editor.Resize(size)
						a.editor.Refresh()
					}
				})
			case <-a.done:
				return
			}
		}
	}()
}

// ════════════════════════════════════════════════════════════════
// CURSOR BLINK
// ════════════════════════════════════════════════════════════════

func (a *App) startCursorBlink() {
	a.cursorTick = time.NewTicker(530 * time.Millisecond)

	go func() {
		for {
			select {
			case <-a.cursorTick.C:
				fyne.Do(func() {
					a.editor.ToggleCursor()
				})
			case <-a.done:
				return
			}
		}
	}()
}

// ════════════════════════════════════════════════════════════════
// ACCESSORS
// ════════════════════════════════════════════════════════════════

// Window returns the main window.
func (a *App) Window() fyne.Window {
	return a.window
}

// Editor returns the editor widget.
func (a *App) Editor() *EditorWidget {
	return a.editor
}
