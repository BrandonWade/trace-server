package contact

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Connection - represents a websocket connection
type Connection struct {
	conn     *websocket.Conn
	buffSize int
}

// NewConnection - return a new Connection
func NewConnection(buffSize int) *Connection {
	return &Connection{
		buffSize: buffSize,
	}
}

// Open - upgrades an http request into a websocket connection
func (c *Connection) Open(w *http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  c.buffSize,
		WriteBufferSize: c.buffSize,
	}

	c.conn, err := upgrader.Upgrade(*w, r, nil)
	if err != nil {
		log.Fatal("error opening websocket", err)
	}
}

// Close - closes an open websocket connection
func (c *Connection) Close() {
	c.conn.WriteMessage(websocket.CloseNormalClosure, []byte(""))
	c.conn.Close()
}
