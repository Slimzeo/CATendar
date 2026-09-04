package main

// Wails owns the macOS AppDelegate, which conflicts with getlantern/systray.
// Keep the normal Dock lifecycle on macOS and retain the tray on other targets.
func startSystray() {}

func stopSystray() {}

func hideWindowOnClose() bool {
	return false
}
