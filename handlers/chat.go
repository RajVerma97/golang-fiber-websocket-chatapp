// handlers/chat.go
package handlers

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Message struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	Username string `json:"username"`
}

type Client struct {
	Conn     *websocket.Conn
	Username string
}

var (
	clients   = make(map[*websocket.Conn]*Client)
	broadcast = make(chan Message)
	mutex     = sync.RWMutex{}
)

func ChatHandler(c *websocket.Conn) {
	// Create new client
	client := &Client{
		Conn: c,
	}

	// Register client
	mutex.Lock()
	clients[c] = client
	mutex.Unlock()

	log.Printf("New client connected: %s\n", c.RemoteAddr())

	// Cleanup on disconnect
	defer func() {
		mutex.Lock()
		delete(clients, c)
		mutex.Unlock()
		log.Printf("Client disconnected: %s\n", c.RemoteAddr())
		c.Close()
	}()

	// Send welcome message
	welcomeMsg := Message{
		Type:     "system",
		Content:  "Welcome to the chat!",
		Username: "System",
	}
	if err := c.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Error sending welcome message: %v\n", err)
		return
	}

	// Message handling loop
	for {
		var msg Message
		if err := c.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v\n", err)
			}
			return
		}

		// Update client username if not set
		if clients[c].Username == "" {
			clients[c].Username = msg.Username
		}

		// Log and broadcast message
		log.Printf("Message received from %s: %s\n", msg.Username, msg.Content)
		broadcast <- msg
	}
}

func init() {
	// Start broadcast handler
	go handleBroadcast()
}

func handleBroadcast() {
	for {
		msg := <-broadcast
		mutex.RLock()
		for client := range clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("Error broadcasting message: %v\n", err)
				client.Close()
				mutex.RUnlock()
				mutex.Lock()
				delete(clients, client)
				mutex.Unlock()
				mutex.RLock()
			}
		}
		mutex.RUnlock()
	}
}
