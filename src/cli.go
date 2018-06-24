package yolo

import (
	"bytes"
	"os"
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

	if err := cmd.Run(); err != nil {
		return stdout.String(), stderr.String(), err
	}

	log.Info("Done '%s'", command)

	return stdout.String(), stderr.String(), nil
}

func CurrentGitBranch() string {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return ""
	}

	return stdout.String()
}

func WorkingDir() string {
	wd, _ := os.Getwd()
	return wd
}
