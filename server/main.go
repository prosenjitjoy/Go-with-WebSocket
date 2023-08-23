package main

import (
	"context"
	"fmt"
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

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

		var v interface{}
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			log.Fatal("Unable to read messages:", err)
		}

		msg := fmt.Sprintf("Received: %+v\n", v)
		err = wsjson.Write(ctx, c, msg)
		if err != nil {
			log.Fatal("Unable to write messages:", err)
		}

		fmt.Println(msg)
		c.Close(websocket.StatusNormalClosure, "")
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
