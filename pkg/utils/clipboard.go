package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// CopyToClipboard copies text to the system clipboard
func CopyToClipboard(content string) error {
	switch runtime.GOOS {
	case "darwin":
		return copyToClipboardMac(content)
	case "linux":
		return copyToClipboardLinux(content)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// copyToClipboardMac copies text to clipboard on macOS
func copyToClipboardMac(content string) error {
	// Try pbcopy first (most common)
	_, err := exec.LookPath("pbcopy")
	if err == nil {
		return runCommand("pbcopy", content)
	}

	// Fallback to xclip/xsel if available (for macOS with X11)
	_, err = exec.LookPath("xclip")
	if err == nil {
		return runCommand("xclip", "-selection", "clipboard", content)
	}

	_, err = exec.LookPath("xsel")
	if err == nil {
		return runCommand("xsel", "--clipboard", "--input", content)
	}

	return fmt.Errorf("no clipboard utility found (tried: pbcopy, xclip, xsel)")
}

// copyToClipboardLinux copies text to clipboard on Linux
func copyToClipboardLinux(content string) error {
	// Try xclip first
	_, err := exec.LookPath("xclip")
	if err == nil {
		return runCommand("xclip", content, "-selection", "clipboard")
	}

	// Try xsel as alternative
	 _, err = exec.LookPath("xsel")
	if err == nil {
		return runCommand("xsel", content, "--clipboard", "--input")
	}

	// Try wl-copy for Wayland
	_, err = exec.LookPath("wl-copy");
	if err == nil {
		return runCommand("wl-copy", content)
	}

	// Try termux-api for Android/Termux
	 _, err = exec.LookPath("termux-clipboard-set")
	if err == nil {
		return runCommand("termux-clipboard-set", content)
	}

	return fmt.Errorf("no clipboard utility found (tried: xclip, xsel, wl-copy, termux-clipboard-set)")
}

// runCommand runs a command with the given content as input
func runCommand(name string, content string, args ...string) error {
	cmd := exec.Command(name, args...)

	// Get stdin pipe
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Write content to stdin
	_, err = stdin.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("failed to write to stdin: %w", err)
	}

	// Close stdin
	err = stdin.Close()
	if err != nil {
		return fmt.Errorf("failed to close stdin: %w", err)
	}

	// Wait for command to complete
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

// IsClipboardAvailable checks if clipboard functionality is available
func IsClipboardAvailable() bool {
	switch runtime.GOOS {
	case "darwin":
		return checkCommandExists("pbcopy") || checkCommandExists("xclip") || checkCommandExists("xsel")
	case "linux":
		return checkCommandExists("xclip") || checkCommandExists("xsel") || checkCommandExists("wl-copy") || checkCommandExists("termux-clipboard-set")
	case "windows":
		return checkCommandExists("powershell") || checkCommandExists("clip")
	default:
		return false
	}
}

// checkCommandExists checks if a command exists in PATH
func checkCommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// GetClipboardInfo returns information about available clipboard utilities
func GetClipboardInfo() string {
	var available []string

	switch runtime.GOOS {
	case "darwin":
		if checkCommandExists("pbcopy") {
			available = append(available, "pbcopy")
		}
		if checkCommandExists("xclip") {
			available = append(available, "xclip")
		}
		if checkCommandExists("xsel") {
			available = append(available, "xsel")
		}
	case "linux":
		if checkCommandExists("xclip") {
			available = append(available, "xclip")
		}
		if checkCommandExists("xsel") {
			available = append(available, "xsel")
		}
		if checkCommandExists("wl-copy") {
			available = append(available, "wl-copy")
		}
		if checkCommandExists("termux-clipboard-set") {
			available = append(available, "termux-clipboard-set")
		}
	case "windows":
		if checkCommandExists("powershell") {
			available = append(available, "PowerShell")
		}
		if checkCommandExists("clip") {
			available = append(available, "clip")
		}
	}

	if len(available) == 0 {
		return "No clipboard utilities found"
	}

	return fmt.Sprintf("Available: %s", available)
}
