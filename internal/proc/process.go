package proc

import (
	"io"
	"os/exec"
)

// ExecCommand is a variable so tests can override the command construction.
var ExecCommand = exec.Command

// Run starts the provided command and optionally waits for it to exit.
// On Start error it returns nil after logging should be done by the caller; on
// Wait error it returns nil as well (maintains previous behavior of not
// propagating process errors).
func Run(cmd *exec.Cmd, shouldWait bool, stderr io.Writer) error {
	if err := cmd.Start(); err != nil {
		if stderr != nil {
			name := cmd.Path
			if name == "" && len(cmd.Args) > 0 {
				name = cmd.Args[0]
			}
			_, _ = stderr.Write([]byte("Failed to start process '" + name + "': " + err.Error() + "\n"))
		}
		return nil
	}
	if shouldWait {
		if err := cmd.Wait(); err != nil {
			if stderr != nil {
				name := cmd.Path
				if name == "" && len(cmd.Args) > 0 {
					name = cmd.Args[0]
				}
				_, _ = stderr.Write([]byte("Process '" + name + "' exited with error: " + err.Error() + "\n"))
			}
		}
	}
	return nil
}
