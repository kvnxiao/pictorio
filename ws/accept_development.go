// +build development

package ws

import (
	"net/http"

	"nhooyr.io/websocket"
)

func Accept(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:8080"},
	})
}
