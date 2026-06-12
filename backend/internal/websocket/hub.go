// Copyright 2026. Kimjibeom. All rights reserved.
package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

// EventType defines the type of WebSocket event.
type EventType string

const (
	EventNewReservation    EventType = "NEW_RESERVATION"
	EventReservationUpdate EventType = "RESERVATION_UPDATE"
	EventWaitingQueueUpdate EventType = "WAITING_QUEUE_UPDATE"
	EventCIDIncoming       EventType = "CID_INCOMING"
	EventOnlineBooking     EventType = "ONLINE_BOOKING"
	EventNotification      EventType = "NOTIFICATION"
)

// Event represents a WebSocket message to broadcast.
type Event struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
	Time    string      `json:"time"`
}

// Client represents a single WebSocket connection.
type Client struct {
	hub  *Hub
	conn *ws.Conn
	send chan []byte
}

// Hub manages all WebSocket clients and broadcasts events.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("WebSocket client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastEvent sends an event to all connected clients.
func (h *Hub) BroadcastEvent(eventType EventType, payload interface{}) {
	event := Event{
		Type:    eventType,
		Payload: payload,
		Time:    time.Now().Format(time.RFC3339),
	}
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal WebSocket event: %v", err)
		return
	}
	h.broadcast <- data
}

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO(security): In production, validate origin against allowed origins list
		origin := r.Header.Get("Origin")
		// For development, allow localhost origins
		return origin == "http://localhost:3000" || origin == "https://localhost:3000" || origin == ""
	},
}

// HandleWebSocket is the Gin handler for WebSocket connections.
func (h *Hub) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}
	h.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(ws.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
