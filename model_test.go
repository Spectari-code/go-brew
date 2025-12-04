package main

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TestInitialModel verifies that the initial model is created with the correct
// default values and configuration. This ensures the application starts in
// a predictable state with proper timer values and state initialization.
func TestInitialModel(t *testing.T) {
	config := NewConfig() // Use NewConfig to get presets
	config.BrewTime = 5 * time.Minute
	model := initialModel(config)

	if model.timer != 5*time.Minute {
		t.Errorf("Expected timer %v, got %v", 5*time.Minute, model.timer)
	}
	if model.config.BrewTime != 5*time.Minute {
		t.Errorf("Expected config BrewTime %v, got %v", 5*time.Minute, model.config.BrewTime)
	}
	if model.state != StateIdle {
		t.Errorf("Expected state %v, got %v", StateIdle, model.state)
	}
	if model.presetIdx != 0 {
		t.Errorf("Expected presetIdx %v, got %v", 0, model.presetIdx)
	}
}

// TestUpdateStart verifies that pressing the start key ('s') transitions the model
// from idle state to brewing state and initiates the timer tick command.
func TestUpdateStart(t *testing.T) {
	config := NewConfig() // Use NewConfig to get presets
	config.BrewTime = 1 * time.Minute
	mdl := initialModel(config)

	newModel, cmd := mdl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	m, ok := newModel.(model)
	if !ok {
		t.Fatal("Failed to cast to model")
	}

	if !m.isBrewing() {
		t.Error("Expected brewing to be true")
	}
	if cmd == nil {
		t.Error("Expected cmd to be not nil")
	}
}

// TestUpdateReset verifies that pressing the reset key ('r') transitions the model
// from any state back to idle and resets the timer to the current preset duration.
func TestUpdateReset(t *testing.T) {
	config := NewConfig() // Use NewConfig to get presets
	config.BrewTime = 1 * time.Minute
	mdl := initialModel(config)
	mdl.state = StateBrewing
	mdl.timer = 30 * time.Second

	newModel, _ := mdl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	m, ok := newModel.(model)
	if !ok {
		t.Fatal("Failed to cast to model")
	}

	if m.isBrewing() {
		t.Error("Expected brewing to be false")
	}
	// Reset uses preset duration, not config.BrewTime
	expectedDuration := m.currentPreset().Duration
	if m.timer != expectedDuration {
		t.Errorf("Expected timer %v, got %v", expectedDuration, m.timer)
	}
}

// TestView verifies that the View function generates a non-empty string containing
// expected UI elements for the idle state. This ensures the UI renders correctly.
func TestView(t *testing.T) {
	config := NewConfig() // Use NewConfig to get presets
	config.BrewTime = 1 * time.Minute
	model := initialModel(config)
	model.width = 80
	model.height = 24

	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}
	// Check if contains expected strings
	if !contains(view, "Press 's' to start") {
		t.Error("Expected start message in view")
	}
}

// TestUpdatePauseResume verifies that the spacebar key correctly toggles between
// brewing and paused states, demonstrating proper state machine transitions.
func TestUpdatePauseResume(t *testing.T) {
	config := NewConfig()
	config.BrewTime = 1 * time.Minute
	mdl := initialModel(config)

	// Start brewing
	newModel, _ := mdl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	m, ok := newModel.(model)
	if !ok {
		t.Fatal("Failed to cast to model")
	}

	if !m.isBrewing() {
		t.Error("Expected brewing to be true after start")
	}

	// Press spacebar to pause
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = newModel.(model)

	if !m.isPaused() {
		t.Error("Expected paused to be true after spacebar")
	}

	// Press spacebar again to resume
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = newModel.(model)

	if !m.isBrewing() {
		t.Error("Expected brewing to be true after spacebar resume")
	}
}

// TestCustomDurationPrecedence verifies that when a custom duration is set via command line,
// it takes precedence over tea preset durations when starting the timer.
func TestCustomDurationPrecedence(t *testing.T) {
	config := NewConfig()
	config.BrewTime = 2 * time.Minute  // Custom duration
	config.CustomDuration = true       // Simulate -duration flag being used
	mdl := initialModel(config)

	// Start timer
	newModel, _ := mdl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	m, ok := newModel.(model)
	if !ok {
		t.Fatal("Failed to cast to model")
	}

	// Verify timer remains at custom duration (2 minutes), not preset duration
	if m.timer != 2*time.Minute {
		t.Errorf("Expected timer %v, got %v", 2*time.Minute, m.timer)
	}
}

// TestCustomDurationReset verifies that custom duration is preserved when resetting timer.
func TestCustomDurationReset(t *testing.T) {
	config := NewConfig()
	config.BrewTime = 3 * time.Minute  // Custom duration
	config.CustomDuration = true       // Simulate -duration flag being used
	mdl := initialModel(config)

	// Start timer
	mdl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})

	// Reset timer
	newModel, _ := mdl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	m, ok := newModel.(model)
	if !ok {
		t.Fatal("Failed to cast to model")
	}

	// Verify custom duration is preserved after reset
	if m.timer != 3*time.Minute {
		t.Errorf("Expected timer %v, got %v", 3*time.Minute, m.timer)
	}
}

// TestPresetNavigationWithCustomDuration verifies that when custom duration is set,
// navigating through presets doesn't change the timer duration.
func TestPresetNavigationWithCustomDuration(t *testing.T) {
	config := NewConfig()
	config.BrewTime = 5 * time.Minute  // Custom duration
	config.CustomDuration = true       // Simulate -duration flag being used
	mdl := initialModel(config)

	// Navigate through presets
	newModel, _ := mdl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("up")})
	mdl, _ = newModel.(model)
	newModel, _ = mdl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("down")})
	m, _ := newModel.(model)

	// Verify custom duration is preserved during preset navigation
	if m.timer != 5*time.Minute {
		t.Errorf("Expected timer %v, got %v", 5*time.Minute, m.timer)
	}
}

// TestDefaultBehaviorWithoutCustomDuration verifies that when no custom duration is set,
// the application behaves as before using preset durations.
func TestDefaultBehaviorWithoutCustomDuration(t *testing.T) {
	config := NewConfig()
	config.BrewTime = DefaultBrewTime     // Use default
	config.CustomDuration = false         // No custom duration
	mdl := initialModel(config)

	// Start timer
	newModel, _ := mdl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	m, ok := newModel.(model)
	if !ok {
		t.Fatal("Failed to cast to model")
	}

	// Verify timer uses preset duration, not config.BrewTime
	expectedDuration := m.currentPreset().Duration
	if m.timer != expectedDuration {
		t.Errorf("Expected timer %v (preset), got %v", expectedDuration, m.timer)
	}
}

// TestPresetNavigationWithoutCustomDuration verifies that preset navigation works
// normally when no custom duration is set.
func TestPresetNavigationWithoutCustomDuration(t *testing.T) {
	config := NewConfig()
	config.BrewTime = DefaultBrewTime     // Use default
	config.CustomDuration = false         // No custom duration
	mdl := initialModel(config)
	originalPresetIdx := mdl.presetIdx

	// Navigate to next preset
	newModel, _ := mdl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("down")})
	m, ok := newModel.(model)
	if !ok {
		t.Fatal("Failed to cast to model")
	}

	// Verify preset changed and timer updated to new preset duration
	if m.presetIdx == originalPresetIdx {
		t.Error("Expected preset index to change")
	}
	expectedDuration := m.currentPreset().Duration
	if m.timer != expectedDuration {
		t.Errorf("Expected timer %v (new preset), got %v", expectedDuration, m.timer)
	}
}

// contains is a helper function that checks if a substring exists within a string.
// It uses a recursive approach for substring searching without relying on strings.Contains.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))
}
