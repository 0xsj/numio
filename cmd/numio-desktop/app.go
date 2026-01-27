package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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
	a.window.Resize(fyne.NewSize(800, 600))
	a.window.CenterOnScreen()

	// Create editor widget
	a.editor = NewEditorWidget()

	// Set up layout
	content := container.NewMax(a.editor)
	a.window.SetContent(content)

	// Focus the editor
	a.window.Canvas().Focus(a.editor)

	// Start cursor blink
	a.startCursorBlink()

	// Handle window close
	a.window.SetOnClosed(func() {
		a.Stop()
	})

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
// CURSOR BLINK
// ════════════════════════════════════════════════════════════════

func (a *App) startCursorBlink() {
	a.cursorTick = time.NewTicker(500 * time.Millisecond)

	go func() {
		for {
			select {
			case <-a.cursorTick.C:
				a.editor.ToggleCursor()
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
