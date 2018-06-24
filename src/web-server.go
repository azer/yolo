package yolo

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var (
	open     = []*websocket.Conn{}
	rw       sync.RWMutex
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func CreateWebSocket(onMessage func(string)) func(http.ResponseWriter, *http.Request) {
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

			go onMessage(string(msg))

		}

		rw.Lock()
		open = append(open[:index], open[index+1:]...)
		rw.Unlock()

		log.Info("Client %d disconnected", index)
	}
}

func SendMessage(content []byte) error {
	log.Info("Send message to browser %s", string(content))

	if len(content) == 0 {
		log.Info("No open connections")
	}

	for _, conn := range open {
		if err := conn.WriteMessage(websocket.TextMessage, content); err != nil {
			return err
		}
	}

	return nil
}

func WebServer(addr string, webInterface func(http.ResponseWriter, *http.Request), onMessage func(string)) {
	http.HandleFunc("/socket", CreateWebSocket(onMessage))
	http.HandleFunc("/", webInterface)
	http.ListenAndServe(addr, nil)
}
