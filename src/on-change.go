package yolo

import (
	"time"
)

var timer *time.Timer

func RunOnChange(build *Build) func(*WatchEvent) {
	return func(event *WatchEvent) {
		if timer != nil {
			timer.Stop()
			timer = nil
		}

		timer = time.NewTimer(time.Millisecond * 2000)

		go DebounceExecution(build)
	}
}

func DebounceExecution(build *Build) {
	if (build.Started && !build.Done) {
		log.Info("Event ignored due to executing build")
		return
	}
	
	if timer != nil {
		<-timer.C
		timer = nil
	}

	build.Reset()
	build.Started = true

	build.Distribute()
	build.Execute()
	build.Distribute()
}
