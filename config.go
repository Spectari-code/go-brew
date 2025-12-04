package main

import (
	"flag"
	"fmt"
	"time"
)

// Constants contain application-wide configuration values and defaults.
const (
	DefaultBrewTime      = 4 * time.Minute
	MinBrewTime         = 30 * time.Second
	MaxBrewTime         = 30 * time.Minute
	DefaultProgressBarWidth = 20

	// Colors
	ColorReady   = "#00FF7F"
	ColorBrewing = "#FFD93D"
	ColorPaused  = "#FFA500"
	ColorIdle    = "#AAAAAA"

	// Keys
	KeyStart   = "s"
	KeyReset   = "r"
	KeyQuit    = "q"
	KeyQuitAlt = "ctrl+c"
	KeyPause   = "space"
	KeyUp      = "up"
	KeyDown    = "down"
)

// TimerState represents the current state of the timer in the brewing lifecycle.
// It follows a simple state machine pattern for predictable behavior.
type TimerState int

const (
	// StateIdle represents the initial state where the timer is not running
	StateIdle TimerState = iota
	// StateBrewing indicates the timer is actively counting down
	StateBrewing
	// StatePaused indicates the timer is temporarily stopped and can be resumed
	StatePaused
	// StateFinished indicates the timer has completed and tea is ready
	StateFinished
)

// KeyBinding represents a keyboard shortcut and its user-facing description.
// This provides a flexible way to map keyboard input to actions.
type KeyBinding struct {
	Key  string // The keyboard key or combination
	Desc string // Human-readable description of the action
}

// TeaPreset represents a pre-configured tea brewing setting with all necessary
// information for proper tea preparation. Each preset includes brew time,
// recommended temperature, and helpful notes for the best results.
type TeaPreset struct {
	Name     string        // Human-readable name of the tea type
	Duration time.Duration // Recommended brewing time
	Temp     string        // Recommended water temperature
	Notes    string        // Additional brewing notes or tips
}

// DefaultTeaPresets contains carefully selected tea presets for common tea types.
// These presets are based on standard brewing recommendations and provide
// excellent starting points for different tea varieties.
var DefaultTeaPresets = []TeaPreset{
	{"Rooibos", 4 * time.Minute, "95°C", "No bitterness, naturally sweet"},
	{"Green Tea", 2 * time.Minute, "80°C", "Don't overbrew to avoid bitterness"},
	{"Black Tea", 3 * time.Minute, "95°C", "Full flavor development"},
	{"Herbal", 5 * time.Minute, "95°C", "Medicinal properties develop over time"},
	{"White Tea", 2 * time.Minute, "75°C", "Delicate flavor, careful timing"},
	{"Oolong", 3 * time.Minute, "85°C", "Complex flavors, multiple infusions possible"},
}

// Config holds all application configuration including user settings,
// tea presets, key bindings, and preferences. It provides a centralized
// location for all configurable aspects of the application.
type Config struct {
	BrewTime      time.Duration // Default brew time when no preset is selected
	SoundEnabled  bool          // Whether to play audio alerts when tea is ready
	NotifyEnabled bool          // Whether to show desktop notifications
	KeyBindings   []KeyBinding  // List of keyboard shortcuts and their descriptions
	Presets       []TeaPreset   // Available tea presets with their brewing parameters
}

// NewConfig creates a new Config instance with sensible default values.
// The configuration includes all standard tea presets, default key bindings,
// and enabled audio/notification features for the best user experience.
func NewConfig() *Config {
	return &Config{
		BrewTime:      DefaultBrewTime,
		SoundEnabled:  true,
		NotifyEnabled: true,
		Presets:       DefaultTeaPresets,
		KeyBindings: []KeyBinding{
			{"s", "Start timer"},
			{KeyPause, "Pause/Resume"},
			{"r", "Reset timer"},
			{KeyUp + "/" + KeyDown, "Select preset"},
			{"q/ctrl+c", "Quit"},
		},
	}
}

// Validate checks that the configuration values are within acceptable ranges.
// This prevents invalid configurations that could cause runtime errors or
// poor user experience. Returns an error if validation fails.
func (c *Config) Validate() error {
	if c.BrewTime < MinBrewTime {
		return fmt.Errorf("brew time must be at least %v", MinBrewTime)
	}
	if c.BrewTime > MaxBrewTime {
		return fmt.Errorf("brew time cannot exceed %v", MaxBrewTime)
	}
	return nil
}

// ParseFlags parses command line flags and updates the configuration accordingly.
// Currently supports the -duration flag for custom brew times.
// This should be called after NewConfig() but before Validate().
func (c *Config) ParseFlags() {
	flag.DurationVar(&c.BrewTime, "duration", c.BrewTime, "brew time for the tea timer")
	flag.Parse()
}
