package tsgen

import (
	"os/exec"
)

func formatCode(path string, formatter string, args ...string) error {
	args = append(args, path)
	cmd := exec.Command(formatter, args...)
	return cmd.Run()
}
