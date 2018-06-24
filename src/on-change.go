package yolo

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var timer *time.Timer

func RunOnChange(command string) func(*WatchEvent) {
	return func(event *WatchEvent) {
		if timer != nil {
			timer.Stop()
			timer = nil
		}

		timer = time.NewTimer(time.Millisecond * 2000)

		go DebounceExecution(command)
	}
}

func DebounceExecution(command string) {
	if timer != nil {
		<-timer.C
		timer = nil
	}

	msg := &struct {
		Started bool   `json:"started"`
		Done    bool   `json:"done"`
		Command string `json:"command"`
		Stdout  string `json:"stdout"`
		Stderr  string `json:"stderr"`
	}{true, false, command, "", ""}

	started, _ := json.Marshal(msg)

	SendMessage(started)

	log.Info("exec ")
	stdout, stderr, err := ExecuteCommand(command)
	if err != nil {
		fmt.Fprintf(os.Stderr, stderr)
	}

	log.Info("!exec %s", stderr)

	msg.Done = true
	msg.Stdout = stdout
	msg.Stderr = stderr

	done, _ := json.Marshal(msg)

	SendMessage(done)
}
