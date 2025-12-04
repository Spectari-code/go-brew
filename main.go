// Package main implements a CLI tea timer application using the Bubbletea TUI framework.
//
// The application follows the Model-View-Update (MVU) architecture pattern:
//   - Model: Application state including timer, brewing status, and configuration
//   - View: Renders the terminal UI with centered display, progress bar, and controls
//   - Update: Handles user input and timer updates, managing state transitions
//
// Features include:
//   - Configurable brew times with built-in tea presets
//   - Visual progress bar with color-coded states
//   - Cross-platform audio notifications and desktop alerts
//   - Keyboard controls for start, pause, reset, and quit operations
//   - Responsive design that adapts to terminal size
//
// Usage:
//   go run .                    # Run with default settings
//   go run . -duration 2m       # Run with 2-minute timer
//
// Key controls:
//   s, space     - Start/pause timer
//   r            - Reset timer
//   up/down      - Select tea preset
//   q, ctrl+c    - Quit application
package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

// These variables are set by GoReleaser during build
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "local"
)

// Init initializes the Bubbletea program with no initial commands.
// This is called once when the program starts and sets up the initial state.
func (m model) Init() tea.Cmd {
	return nil
}

// printVersion prints version information and exits
func printVersion() {
	fmt.Printf("go-brew %s\n", version)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Built: %s\n", date)
	fmt.Printf("Built by: %s\n", builtBy)
}

// main is the entry point of the Go Brew CLI application.
// It sets up the configuration, validates it, and starts the Bubbletea TUI program.
// The program runs in alternate screen mode for a full terminal experience.
func main() {
	config := NewConfig()
	config.ParseFlags()

	// Handle version flag
	if config.ShowVersion {
		printVersion()
		return
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	p := tea.NewProgram(initialModel(config), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Printf("Error running program: %v", err)
	}
}
