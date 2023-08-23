package main

import (
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{"localhost", "example.com"},
		})
		if err != nil {
			log.Fatal("Unable to establish handshaking:", err)
		}
		defer c.Close(websocket.StatusInternalError, "the sky is falling")
		c.Close(websocket.StatusNormalClosure, "cross origin WebSocket accepted")
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
