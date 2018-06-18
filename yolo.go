package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/azer/yolo/src"
	"net/http"
	"text/template"
)

func main() {
	var (
		include yolo.Patterns
		exclude yolo.Patterns
	)

	flag.Var(&include, "i", "Glob pattern to watch for changes")
	flag.Var(&exclude, "e", "Glob pattern to ignore")
	command := flag.String("c", "", "Command to execute on change")

	flag.Parse()

	watch, err := yolo.NewWatch(&include, &exclude)
	if err != nil {
		panic(err)
	}

	go yolo.WebServer(":8080", WebInterface, OnMessage)

	watch.Start(func(event *yolo.WatchEvent) {
		msg := &struct {
			Started bool   `json:"started"`
			Done    bool   `json:"done"`
			Command string `json:"command"`
			Stdout  string `json:"stdout"`
			Stderr  string `json:"stderr"`
		}{true, false, *command, "", ""}

		started, _ := json.Marshal(msg)

		yolo.SendMessage(started)

		stdout, stderr, err := yolo.ExecuteCommand(*command)

		msg.Done = true
		msg.Stdout = stdout
		msg.Stderr = stderr

		done, _ := json.Marshal(msg)

		yolo.SendMessage(done)

		if err != nil {
			fmt.Println(stderr)
		}
	})
}

func OnMessage(msg string) {
	fmt.Println("Received %v", string(msg))
}

func WebInterface(w http.ResponseWriter, r *http.Request) {
	homepage := template.Must(template.New("homepage").Parse(html))

	homepage.Execute(w, map[string]string{
		"CSS":  css,
		"ADDR": "localhost:8080",
	})
}

const html = `<!DOCTYPE html>
<html>
<head>
  <title>Ready - Yolo</title>
  <style type="text/css">
  {{.CSS}}
  </style>
</head>
<body>
  <div class="error-container container">
    <h1>Error</h1>
    <h2>Build failed due to following errors below:</h2>
    <div class="stderr">
    </div>
  </div>
  <div class="ready-container container">
    <h1>Ready</h1>
    <h2>Listening for changes...</h2>
  </div>
  <div class="busy-container container">
    <h1>Building</h1>
    <h2>$ <span class="command"></span></h2>
  </div>
  <div class="success-container container">
    <h1>Done</h1>
    <h2>Built completed successfully.</h2>
  </div>
  <div class="connection"></div>
  <script type="text/javascript">
var socket = new WebSocket("ws://{{.ADDR}}/socket");

socket.onopen = function () {
  document.querySelector('.connection').innerHTML = "Connected";
};

socket.onclose = function () {
  document.querySelector('.connection').innerHTML = "";
};

socket.onmessage = function (e) {
  const parsed = JSON.parse(e.data)

  console.log('Message received', parsed)

  if (parsed.started) {
    document.title = 'Building...'
    document.body.className = "busy"
    document.querySelector(".command").innerHTML = parsed.command;
  }

  if (parsed.done && parsed.stderr) {
    document.body.className = "error"
    document.querySelector('.stderr').innerHTML = parsed.stderr.split('\n').join('<br />')
  } else if (parsed.done) {
    document.body.className = "success"
  }
};

function send(msg) {
  socket.send(msg);
}
</script>
</body>
</html>
`

const css = `
html, body {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
  font: 400 1rem "-apple-system", "BlinkMacSystemTypography", "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", "sans-serif";
  line-height: 1.6em;
}

body {
  display: flex;
  justify-content: center;
  align-items: center;
  background: #eec660;
  color: #fff;
}

.container {
  max-width: 90%;
  margin: 0 auto;
}

.error-container {
  background: #ff4136;
  color: #fbf1a9;
}

h1 {
  font-size: 6rem;
  margin-bottom: 2.2rem;
  text-transform: uppercase;
}

h2 {
  font-size: 1.3rem;
  font-weight: 400;
  margin-left: 7px;
}

.stderr {
  margin: 10px 0;
  font: 400 1rem "Menlo", "Inconsolata", "Fira Mono", "Noto Mono", "Droid Sans Mono", "Consolas", "monaco" , "monospace";
  line-height: 2em;
  color: #fff;
  padding: 25px;
  background: rgba(0, 0, 0, 0.1);
  list-style: none;
  box-sizing: border-box;
}

.error-container, .busy-container, .success-container {
  display: none;
}

.connection {
  position: absolute;
  bottom: 20px;
  right: 20px;
  color: rgba(255, 255, 255, 0.5);
  text-transform: uppercase;
  font: 14px "Menlo", "Inconsolata", "Fira Mono", "Noto Mono", "Droid Sans Mono", "Consolas", "monaco" , "monospace";
}

body.busy .ready-container, body.success .ready-container, body.error .ready-container {
  display: none;
}

body.busy .busy-container, body.success .success-container, body.error .error-container {
  display: block;
}

body.busy {
  background: #70d9e1;
  color: #fff;
}

body.success {
  background: #13b66a;
  color: #fff;
}

body.error {
  background: #ff4136;
  color: #fbf1a9;
}
`
