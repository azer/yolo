package yolo

import (
	"bytes"
	"os/exec"
)

func ExecuteCommand(command string) (string, string, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	log.Info("Executing %s", command)

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Info("Done '%s'", command)

	if err := cmd.Run(); err != nil {
		return stdout.String(), stderr.String(), err
	}

	return stdout.String(), stderr.String(), nil
}
