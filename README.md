# yolo

A file watcher with web based user interface, so you can see the status of your build (and error output, if any) instantly.

![](https://cldup.com/G0VmmMWMnz.gif)

Usage example:

```bash
$ yolo -i *.go -c 'go build' -a localhost:8080
```

It's silent by default. You can enable internal logging by;

```bash
$ LOG=* yolo ...
```

# Install

```bash
$ go get github.com/azer/yolo
```

# Todo

- [ ] Read filenames from another file. So we could do `yolo -f index.html` and get live updates when dependencies change.
- [ ] Show last build status
- [ ] Escape characters
- [ ] Add option for reading error messages from stdout
- [ ] Output the name of changed file
- [ ] Reconnect
- [ ] Split JS to another endpoint so it can be included by other pages
- [ ] How could it be used for viewing web pages / apps ?
- [ ] How could output be processed ?
