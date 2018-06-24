package yolo

import (
	"encoding/json"
	"fmt"
	"os"
)

func NewBuild(command string) *Build {
	return &Build{
		Command:    command,
		WorkingDir: WorkingDir(),
		GitBranch:  CurrentGitBranch(),
	}
}

type Build struct {
	Command    string `json:"command"`
	Started    bool   `json:"started"`
	Done       bool   `json:"done"`
	Failed     bool   `json:"failed"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	WorkingDir string `json:"working_dir"`
	GitBranch  string `json:"git_branch"`
}

func (build *Build) Reset() {
	build.Done = false
	build.Failed = false
	build.Stdout = ""
	build.Stderr = ""
	build.Started = false
}

func (build *Build) Distribute() error {
	msg, err := build.Message()
	if err != nil {
		return err
	}

	return DistributeMessage(msg)
}

func (build *Build) Execute() {
	stdout, stderr, err := ExecuteCommand(build.Command)
	if err != nil {
		fmt.Fprintf(os.Stderr, stderr)
	}

	build.Done = true
	build.Failed = err != nil
	build.Stdout = stdout
	build.Stderr = stderr
}

func (build *Build) Message() ([]byte, error) {
	msg, err := json.Marshal(build)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
