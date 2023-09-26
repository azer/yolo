package main

import (
	"flag"
	"net/http"
	"text/template"

	yolo "github.com/OVYA/yolo/src"
)

func main() {
	var (
		include yolo.Patterns
		exclude yolo.Patterns
	)

	flag.Var(&include, "i", "Glob pattern to include files for watching")
	flag.Var(&exclude, "e", "Glob pattern to exclude from watching")
	command := flag.String("c", "", "Command to execute on change")
	addr := flag.String("a", "", "Host and port to run the web server on")

	flag.Parse()

	if len(include) == 0 || len(*command) == 0 {
		flag.PrintDefaults()
		return
	}

	watch, err := yolo.NewWatch(&include, &exclude)
	if err != nil {
		panic(err)
	}

	build := yolo.NewBuild(*command)

	if len(*addr) > 0 {
		go yolo.WebServer(build, *addr, WebInterface(*addr))
	}

	watch.Start(yolo.RunOnChange(build))
}

func WebInterface(addr string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		homepage := template.Must(template.New("homepage").Parse(html))

		homepage.Execute(w, map[string]string{
			"CSS":  css,
			"ADDR": addr,
		})
	}
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
  <div class="info">
    <div class="wd"></div>
    <div class="git-w">* <span class="git"></span></div>
    <div class="command-w">$ <span class="command"></span></div>
    <div class="connection"></div>
  </div>
  <script type="text/javascript">
var reconnectInterval = 1000
var addr = "{{.ADDR}}"
var socket;

open()

function open () {
  var conn = new WebSocket("ws://" + (addr[0] == ":" ? "localhost" + addr : addr) + "/socket");
  conn.onopen = onOpen.bind(null, conn)
  conn.onclose = onClose
  conn.onmessage = onMessage
  return conn
}

function onOpen (conn) {
  document.querySelector('.connection').innerHTML = "Connected";
  conn.send("ping")
};

function onClose () {
  document.querySelector('.connection').innerHTML = "";
  setTimeout(reconnect, reconnectInterval)
};

function onMessage (e) {
  const parsed = JSON.parse(e.data)

  console.log('Message received', parsed)

  document.querySelectorAll(".command").forEach(el => el.innerHTML = parsed.command)
  document.querySelectorAll(".wd").forEach(el => el.innerHTML = parsed.working_dir)
  document.querySelectorAll(".git").forEach(el => el.innerHTML = parsed.git_branch)

  if (parsed.started) {
    document.title = 'Building... - Yolo'
    document.body.className = "busy"
    document.querySelector(".command").innerHTML = parsed.command;
  }

  if (parsed.done && (parsed.failed || parsed.stderr)) {
    document.title = 'Error - Yolo'
    document.body.className = "error"
     document.querySelector('.stderr').innerHTML = (parsed.stderr || parsed.stdout).split('\n').join('<br />')
  } else if (parsed.done) {
    document.title = 'Done - Yolo'
    document.body.className = "success"
  }
};

function send(msg) {
  socket.send(msg);
}

function reconnect() {
  console.log('Reconnecting...')
  socket = open()
  reconnectInterval += 250
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

.info {
  box-sizing: border-box;
  position: absolute;
  width: 100%;
  bottom: 0;
  background: rgba(0, 0, 0, 0.2);
  color: rgba(255, 255, 255, 0.8);

  display: grid;
  grid-template-columns: repeat(4, 1fr);
  grid-column-gap: 10px;
  justify-items: stretch;
  align-items: stretch;


  text-transform: uppercase;
  font: 14px "Menlo", "Inconsolata", "Fira Mono", "Noto Mono", "Droid Sans Mono", "Consolas", "monaco" , "monospace";
  padding: 20px 0;

  text-align: center;
}

.info div {
  border-right: 1px dotted rgba(255, 255, 255, 0.7);
}

.info div:last-child {
  border-right: 0;
}

.info * {
  text-shadow: 1px 2px rgba(0,0,0, 0.1);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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
