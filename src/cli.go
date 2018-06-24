package yolo

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

func ExecuteCommand(command string) (string, string, error) {
	var (
		stdoutBuf, stderrBuf bytes.Buffer
	)

	cmd := exec.Command("sh", "-c", command)

	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}

	stderrIn, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	if err := cmd.Start(); err != nil {
		return "", "", err
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	if err := cmd.Wait(); err != nil {
		return "", "", err
	}

	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())

	return outStr, errStr, nil

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
