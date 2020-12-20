package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print(fmt.Errorf("main: failed to upgrade: %w", err))
			return
		}

		go func() {
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Print(err)
					break
				}
				_ = conn.WriteMessage(websocket.TextMessage, append([]byte("server: "), msg...))
			}
		}()

	})

	if err := http.ListenAndServe(":9900", mux); err != nil {
		log.Fatal(fmt.Errorf("main: failed to serve http: %w", err))
	}
}
