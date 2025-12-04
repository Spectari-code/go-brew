package main

import "time"

// tickMsg is a Bubbletea message type that represents timer tick events.
// It wraps time.Time to provide timing information for timer updates.
type tickMsg time.Time

// model represents the complete application state for the Go Brew CLI.
// It contains all data needed to render the UI and handle user interactions,
// following the Model-View-Update architecture pattern.
type model struct {
	config    *Config      // Application configuration and settings
	timer     time.Duration // Current remaining time on the timer
	state     TimerState   // Current state of the timer (idle, brewing, paused, finished)
	presetIdx int          // Index of the currently selected tea preset
	width     int          // Terminal width for responsive UI layout
	height    int          // Terminal height for responsive UI layout
}

// initialModel creates a new model instance with the given configuration.
// It initializes the timer to the selected preset duration and sets the
// initial state to idle, ready for user interaction.
func initialModel(config *Config) model {
	return model{
		config:    config,
		timer:     config.BrewTime,
		state:     StateIdle,
		presetIdx: 0,
	}
}

// currentPreset returns the currently selected tea preset from the configuration.
// It includes bounds checking to prevent index out of range errors and
// falls back to the first preset if the selected index is invalid.
func (m model) currentPreset() TeaPreset {
	if m.presetIdx >= 0 && m.presetIdx < len(m.config.Presets) {
		return m.config.Presets[m.presetIdx]
	}
	return m.config.Presets[0]
}

// isBrewing returns true if the timer is currently active and counting down.
// This is a convenience method that checks if the state is StateBrewing.
func (m model) isBrewing() bool {
	return m.state == StateBrewing
}

// isPaused returns true if the timer is currently paused and can be resumed.
// This is a convenience method that checks if the state is StatePaused.
func (m model) isPaused() bool {
	return m.state == StatePaused
}

// isFinished returns true if the timer has completed and tea is ready.
// This is a convenience method that checks if the state is StateFinished.
func (m model) isFinished() bool {
	return m.state == StateFinished
}
