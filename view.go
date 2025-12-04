package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// View renders the complete terminal UI for the Go Brew application.
// It follows the MVU pattern by being a pure function that converts
// the current model state into a string representation for display.
// The view includes the timer display, progress bar, preset information,
// and control hints, all centered in the terminal.
func (m model) View() string {
	// Get current tea preset for display information
	preset := m.currentPreset()

	// Format timer display as MM:SS with leading zeros
	timeStr := fmt.Sprintf("%02d:%02d", int(m.timer.Minutes()), int(m.timer.Seconds())%60)

	// Define reusable styles for consistent UI appearance
	baseStyle := lipgloss.NewStyle().Bold(true).Padding(1, 2)
	presetStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Faint(true)

	// Build comprehensive preset information string
	presetInfo := fmt.Sprintf("%s (%s)", preset.Name, preset.Temp)
	if preset.Notes != "" {
		presetInfo += " - " + preset.Notes
	}

	// Generate status message based on current timer state
	var status string
	switch {
	case m.isFinished():
		// Tea is ready - show completion message with time
		status = baseStyle.Foreground(lipgloss.Color(ColorReady)).Render("🫖 Tea Ready!   " + timeStr)
	case m.isBrewing():
		// Currently brewing - show active status with time
		status = baseStyle.Foreground(lipgloss.Color(ColorBrewing)).Render("⏰ Brewing...   " + timeStr)
	case m.isPaused():
		// Timer paused - show paused status with time
		status = baseStyle.Foreground(lipgloss.Color(ColorPaused)).Render("⏸️ Paused   " + timeStr)
	default:
		// Idle state - show start prompt with time
		status = baseStyle.Foreground(lipgloss.Color(ColorIdle)).Render("Press 's' to start   " + timeStr)
	}

	// Add preset information when idle to help users choose tea type
	if m.state == StateIdle {
		status += "\n" + presetStyle.Render("🍵 "+presetInfo)
	}

	// Generate progress bar for active states (brewing, paused, finished)
	var progress string
	if m.isBrewing() || m.isPaused() || m.isFinished() {
		total := preset.Duration
		elapsed := total - m.timer
		progress = "\n" + renderProgressBar(total, elapsed, DefaultProgressBarWidth, m.state)
	}

	// Build control help section
	controls := "\n\nControls:\n"
	for _, binding := range m.config.KeyBindings {
		controls += fmt.Sprintf("%s: %s\n", binding.Key, binding.Desc)
	}

	// Show current selection details when idle for better UX
	if m.state == StateIdle {
		controls += fmt.Sprintf("\nCurrent: %s (%v)\n", preset.Name, preset.Duration)
	}

	// Combine all UI elements into final display
	ui := status + progress + controls

	// Center the entire UI in the terminal window
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		ui,
	)
}

// renderProgressBar renders a visual progress bar with dynamic styling based on timer state.
// It displays the brewing progress using different characters and colors depending on
// whether the timer is brewing, paused, or finished. The progress bar includes a
// percentage display for precise timing information.
func renderProgressBar(total, elapsed time.Duration, width int, state TimerState) string {
	// Guard against division by zero or invalid total duration
	if total == 0 {
		return ""
	}

	// Calculate progress percentage (clamp between 0 and 1)
	percent := float64(elapsed) / float64(total)
	if percent > 1 {
		percent = 1
	}

	// Determine how many characters should be filled in the progress bar
	filled := int(percent * float64(width))
	bar := ""

	// Select appropriate characters based on timer state for visual feedback
	var fillChar, emptyChar string
	switch state {
	case StateBrewing:
		// Active brewing - use solid fill for completed part
		fillChar, emptyChar = "█", "░"
	case StatePaused:
		// Paused state - use shaded characters to indicate pause
		fillChar, emptyChar = "▓", "▒"
	case StateFinished:
		// Complete - show full bar to indicate completion
		fillChar, emptyChar = "█", "█"
	default:
		// Idle/inactive - use outline characters
		fillChar, emptyChar = "░", "░"
	}

	// Build the progress bar string with appropriate characters
	for i := 0; i < filled; i++ {
		bar += fillChar
	}
	for i := filled; i < width; i++ {
		bar += emptyChar
	}

	// Return formatted progress bar with percentage display
	return fmt.Sprintf("[%s] %.0f%%", bar, percent*100)
}
