package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Fatal("Unable to establish handshaking:", err)
		}
		defer c.Close(websocket.StatusInternalError, "the sky is falling")

		ctx, cancel := context.WithTimeout(r.Context(), time.Minute*10)
		defer cancel()

		ctx = c.CloseRead(ctx)

		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				c.Close(websocket.StatusNormalClosure, "")
				return
			case <-ticker.C:
				err := wsjson.Write(ctx, c, "hi")
				if err != nil {
					log.Fatal("Unable to write into websocket connection:", err)
				}
			}
		}
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
