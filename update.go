package main

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gen2brain/beeep"
)

// Update implements the Bubbletea update function for the Go Brew application.
// It processes incoming messages and updates the model state accordingly.
// This function follows the MVU pattern by returning the updated model and
// any commands that should be executed as side effects.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		// Handle spacebar for pause/resume functionality
		// We check both KeyType and string representation for maximum compatibility
		if msg.Type == tea.KeySpace {
			if m.state == StateBrewing {
				// Pause the timer but keep the current time
				m.state = StatePaused
				return m, nil
			} else if m.state == StatePaused {
				// Resume brewing from the paused state
				m.state = StateBrewing
				return m, tick()
			}
		}

		keyStr := msg.String()
		// Debug: uncomment to see what keys are being pressed
		// log.Printf("Key pressed: %s (Type: %d)", keyStr, msg.Type)

		switch keyStr {
		case KeyQuit, KeyQuitAlt:
			return m, tea.Quit
		case KeyStart:
			// Start timer if not already brewing
			if m.state != StateBrewing {
				// If previously finished, reset to idle before starting fresh
				if m.isFinished() {
					if m.config.CustomDuration {
						m.timer = m.config.BrewTime  // Use custom duration
					} else {
						m.timer = m.currentPreset().Duration  // Use preset duration
					}
					m.state = StateIdle
				}
				// Set timer to custom duration or preset duration and start brewing
				if m.config.CustomDuration {
					m.timer = m.config.BrewTime  // Use custom duration
				} else {
					m.timer = m.currentPreset().Duration  // Use preset duration
				}
				m.state = StateBrewing
				return m, tick() // Start the timer tick mechanism
			}
		case KeyPause:
			// Dedicated pause key (in addition to spacebar)
			if m.state == StateBrewing {
				m.state = StatePaused
				return m, nil
			} else if m.state == StatePaused {
				m.state = StateBrewing
				return m, tick()
			}
		case KeyReset:
			// Reset timer to initial state with custom duration or preset duration
			if m.config.CustomDuration {
				m.timer = m.config.BrewTime  // Use custom duration
			} else {
				m.timer = m.currentPreset().Duration  // Use preset duration
			}
			m.state = StateIdle
			return m, nil
		case KeyUp:
			// Navigate to previous preset (only allowed when idle)
			if m.state == StateIdle {
				// Use modulo arithmetic to wrap around the preset list
				m.presetIdx = (m.presetIdx - 1 + len(m.config.Presets)) % len(m.config.Presets)
				// Only update timer if NOT using custom duration
				if !m.config.CustomDuration {
					m.timer = m.currentPreset().Duration
				}
			}
			return m, nil
		case KeyDown:
			// Navigate to next preset (only allowed when idle)
			if m.state == StateIdle {
				m.presetIdx = (m.presetIdx + 1) % len(m.config.Presets)
				// Only update timer if NOT using custom duration
				if !m.config.CustomDuration {
					m.timer = m.currentPreset().Duration
				}
			}
			return m, nil
		}

	case tickMsg:
		// Handle timer tick events - only process if actively brewing
		if m.state == StateBrewing {
			m.timer -= time.Second
			if m.timer <= 0 {
				// Timer completed - transition to finished state
				m.timer = 0
				m.state = StateFinished
				// Launch asynchronous notifications and sounds
				return m, tea.Cmd(func() tea.Msg {
					go func() {
						// Send desktop notification if enabled
						if m.config.NotifyEnabled {
							if err := beeep.Notify("Go Brew Timer", "Your tea is ready!", ""); err != nil {
								log.Printf("Failed to send notification: %v", err)
							}
						}
						// Play alert sound (includes fallback mechanisms)
						playSound()
					}()
					return nil
				})
			}
			// Continue ticking if not finished
			return m, tick()
		}

	case tea.WindowSizeMsg:
		// Update terminal dimensions for responsive UI layout
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// tick creates a Bubbletea command that generates timer tick messages at one-second intervals.
// This is the core timing mechanism for the application, driving the countdown timer.
// The command continues running until explicitly cancelled by stopping timer operations.
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
