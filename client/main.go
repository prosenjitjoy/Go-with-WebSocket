package main

import (
	"context"
	"log"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:8000/ws", nil)
	if err != nil {
		log.Fatal("Unable to perform handshaking with server:", err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	err = wsjson.Write(ctx, c, "hi")
	if err != nil {
		log.Fatal("Unable to write into websocket connection:", err)
	}

	c.Close(websocket.StatusNormalClosure, "")
}
