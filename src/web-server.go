package yolo

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var (
	rw       sync.RWMutex
	open     = []*websocket.Conn{}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func CreateWebSocket(build *Build) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		rw.Lock()
		index := len(open)
		open = append(open, conn)
		rw.Unlock()

		defer conn.Close()

		log.Info("New connection. Total: %d Index: %d", len(open), index)

		for {
			// Read message from browser
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}

			if string(msg) == "ping" {
				pong, err := build.Message()
				if err != nil {
					panic(err)
				}

				if err := conn.WriteMessage(websocket.TextMessage, pong); err != nil {
					panic(err)
				}
			}

		}

		rw.Lock()
		open = append(open[:index], open[index+1:]...)
		rw.Unlock()

		log.Info("Client %d disconnected", index)
	}
}

func DistributeMessage(content []byte) error {
	log.Info("Send message to browser %s", string(content))

	if len(content) == 0 {
		log.Info("No open connections")
	}

	for _, conn := range open {
		conn.WriteMessage(websocket.TextMessage, content)
	}

	return nil
}

func WebServer(build *Build, addr string, webInterface func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc("/socket", CreateWebSocket(build))
	http.HandleFunc("/", webInterface)
	http.ListenAndServe(addr, nil)
}
