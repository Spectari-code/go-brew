package main

import (
	"bytes"
	_ "embed"
	"io"
	"log"
	"os/exec"
	"runtime"
	"time"

	"github.com/gen2brain/malgo"
	"github.com/hajimehoshi/go-mp3"
)

//go:embed alert.mp3
var alertMP3Data []byte

// playSound attempts to play an audio alert when the timer completes.
// It implements a graceful degradation strategy with multiple fallback options:
// 1. Primary: MP3 playback from embedded alert.mp3 data
// 2. Secondary: System-specific sound files
// 3. Tertiary: Terminal bell character
// This ensures users receive notification even on systems with limited audio capabilities.
func playSound() {
	go func() {
		if err := tryMP3Playback(); err != nil {
			log.Printf("MP3 playback failed: %v", err)
			if err := trySystemBeep(); err != nil {
				log.Printf("System beep failed: %v", err)
				log.Println("All audio methods failed")
			}
		}
	}()
}

// tryMP3Playback attempts to play the embedded MP3 alert file using pure Go libraries.
// It uses go-mp3 for decoding and malgo for cross-platform audio playback.
// This method provides the best audio quality and requires no external files.
func tryMP3Playback() error {
	reader := bytes.NewReader(alertMP3Data)
	decoder, err := mp3.NewDecoder(reader)
	if err != nil {
		return err
	}

	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		log.Printf("MALGO: %s", message)
	})
	if err != nil {
		return err
	}
	defer ctx.Uninit()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatF32
	deviceConfig.Playback.Channels = 2
	deviceConfig.SampleRate = uint32(decoder.SampleRate())

	audioData, err := io.ReadAll(decoder)
	if err != nil {
		return err
	}

	var audioIndex int
	onData := func(outputSamples, inputSamples []byte, frameCount uint32) {
		remaining := len(audioData) - audioIndex
		if remaining <= 0 {
			return
		}

		bytesPerFrame := 4 * 2
		toCopy := int(frameCount) * bytesPerFrame
		if toCopy > remaining {
			toCopy = remaining
		}

		copy(outputSamples, audioData[audioIndex:audioIndex+toCopy])
		audioIndex += toCopy
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: onData,
	})
	if err != nil {
		return err
	}
	defer device.Uninit()

	err = device.Start()
	if err != nil {
		return err
	}

	duration := time.Duration(float64(len(audioData)/(4*2)) / float64(decoder.SampleRate()) * float64(time.Second))
	time.Sleep(duration)

	device.Stop()

	return nil
}

// trySystemBeep attempts to play a system-specific beep sound as a fallback mechanism.
// It uses different methods depending on the operating system to provide the best
// chance of successful audio playback when the MP3 file is unavailable.
func trySystemBeep() error {
	switch runtime.GOOS {
	case "windows":
		return playWindowsBeep()
	case "darwin":
		return playMacBeep()
	case "linux":
		return playLinuxBeep()
	default:
		log.Printf("No system beep implementation for %s", runtime.GOOS)
		return nil
	}
}

// playWindowsBeep plays a system beep sound on Windows using PowerShell.
// It leverages the .NET Media.SoundPlayer class to play the system beep sound.
// This method works on modern Windows systems with PowerShell installed.
func playWindowsBeep() error {
	cmd := exec.Command("powershell", "-c", "(New-Object Media.SoundPlayer 'System.Windows.Media.SystemSounds.Beep.wav').PlaySync();")
	return cmd.Run()
}

// playMacBeep plays a system beep sound on macOS using the afplay command.
// It uses the built-in Ping sound file that's available on all macOS systems.
// This provides a native macOS audio experience without additional dependencies.
func playMacBeep() error {
	cmd := exec.Command("afplay", "/System/Library/Sounds/Ping.aiff")
	return cmd.Run()
}

// playLinuxBeep plays a beep sound on Linux systems with multiple fallback methods.
// Linux audio is highly fragmented, so this function tries several common approaches:
// - paplay (PulseAudio)
// - aplay (ALSA)
// - beep command-line utility
// - Terminal bell character as last resort
func playLinuxBeep() error {
	// Try multiple Linux beep methods
	commands := [][]string{
		{"paplay", "/usr/share/sounds/alsa/Front_Left.wav"},
		{"aplay", "/usr/share/sounds/alsa/Front_Center.wav"},
		{"beep", "-f", "1000", "-l", "200"},
		{"echo", "-e", "\a"},
	}

	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return exec.Command("echo", "-e", "\a").Run()
}
