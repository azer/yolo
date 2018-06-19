# yolo

A file watcher with web based user interface, so you can see the status of your build (and error output, if any) instantly.

![](https://cldup.com/G0VmmMWMnz.gif)

Usage example:

```bash
$ yolo -i *.go -c 'go build' # You can now open up localhost:8080 and watch the build progress.
```

# Install

```bash
$ go get github.com/azer/yolo
```

# Todo

- [ ] Add debouncing
- [ ] Output stdout
- [ ] Add option for reading error messages from stdout
- [ ] Show help when no command is provided
