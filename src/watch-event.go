package yolo

type WatchEvent struct {
	Filename string
	Create   bool
	Write    bool
	Rename   bool
	Remove   bool
}

func (event *WatchEvent) ShouldBeNotified() bool {
	return event.Create || event.Write || event.Rename || event.Remove
}
