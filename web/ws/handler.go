package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan interface{}
}

// HubImpl maintains the set of active connections and broadcasts messages to the
// clients.
type HubImpl struct {
	connections map[*Connection]bool
	broadcast   chan interface{}
	listen      chan interface{}
	register    chan *Connection
	unregister  chan *Connection
	stop        chan struct{}
	command     chan<- *BotCommand
}

// NewHub creates a new Hub object
func NewHub(c chan<- *BotCommand) *HubImpl {
	return &HubImpl{
		broadcast:   make(chan interface{}),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
		stop:        make(chan struct{}),
		command:     c,
	}
}

// Run starts the hub worker loop to dispatch messages inside the websocket hub.
func (h *HubImpl) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				h.command <- &BotCommand{Event: "Disconnect"}
				close(c.send)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		case <-h.stop:
			return
		}
	}
}

// Broadcast a message to all connected clients
func (h *HubImpl) Broadcast(v interface{}) {
	h.broadcast <- v
}

// Stop the run loop
func (h *HubImpl) Stop() {
	close(h.stop)
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writeJSON writes the JSON encoding of v to the connection.
func (c *Connection) writeJSON(v interface{}) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteJSON(v)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case payload, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.writeJSON(payload); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Connection) readPump(h *HubImpl) {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		cm := &BotCommand{}
		err := c.ws.ReadJSON(cm)
		if err != nil {
			log.Print(err)
			break
		}

		h.command <- cm
	}
}

// ServeWs handles websocket requests from the peer.
func (h *HubImpl) ServeWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &Connection{send: make(chan interface{}, 256), ws: ws}
	h.register <- c
	go c.writePump()
	c.readPump(h)
}
