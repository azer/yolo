package yolo

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	conn      *websocket.Conn
	connected = false
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func ListenWebSocket(onMessage func(string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _ = upgrader.Upgrade(w, r, nil)
		connected = true

		for {
			// Read message from browser
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			go onMessage(string(msg))

		}
	}
}

func SendMessage(content []byte) error {
	if !connected {
		return nil
	}

	if err := conn.WriteMessage(websocket.TextMessage, content); err != nil {
		return err
	}

	return nil
}

func WebServer(addr string, webInterface func(http.ResponseWriter, *http.Request), onMessage func(string)) {
	http.HandleFunc("/socket", ListenWebSocket(onMessage))
	http.HandleFunc("/", webInterface)
	http.ListenAndServe(addr, nil)
}
