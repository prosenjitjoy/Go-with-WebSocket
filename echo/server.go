package main

import (
	"context"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

// echoServer is the WebSocket echo server implementation.
// It ensures the client speaks the echo subprotocol and
// only allows one message every 100ms with 10 message burst.
type echoServer struct {
	logf func(f string, v ...interface{})
}

func (s echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		s.logf("%v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	if c.Subprotocol() != "echo" {
		c.Close(websocket.StatusPolicyViolation, "client must speak the echo subprotocol")
		return
	}

	l := rate.NewLimiter(rate.Every(time.Microsecond*100), 10)
	for {
		err = echo(r.Context(), c, l)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			s.logf("failed to echo with %v: %v\n", r.RemoteAddr, err)
			return
		}
	}
}

func echo(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	msgType, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, msgType)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	err = w.Close()
	return err
}
