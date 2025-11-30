package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// GetEditor returns the editor executable name and optional default args.
// Priority:
// 1. $EDITOR (split by whitespace; first token must exist on PATH)
// 2. OS-specific candidates checked in order
// It returns an error if no suitable editor is found on PATH.
func GetEditor() (name string, args []string, err error) {
	if e := os.Getenv("EDITOR"); strings.TrimSpace(e) != "" {
		parts := strings.Fields(e)
		cmd := parts[0]
		if _, lookErr := exec.LookPath(cmd); lookErr == nil {
			if len(parts) > 1 {
				return cmd, parts[1:], nil
			}
			return cmd, nil, nil
		}
		return "", nil, fmt.Errorf("EDITOR specified but not found on PATH: %s", cmd)
	}

	var candidates []string
	switch runtime.GOOS {
	case "darwin":
		candidates = []string{"open", "vim", "nano", "vi"}
	case "windows":
		candidates = []string{"notepad"}
	default: // linux/other
		candidates = []string{"vim", "nano", "vi", "xdg-open"}
	}

	for _, c := range candidates {
		if _, lookErr := exec.LookPath(c); lookErr == nil {
			// For xdg-open/open we don't add extra args; caller will pass the file path.
			return c, nil, nil
		}
	}

	return "", nil, fmt.Errorf("no editor found; set $EDITOR")
}
