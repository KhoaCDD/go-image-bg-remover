package socket

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ClientManager manages WebSocket connections
type ClientManager struct {
	clients map[string]*websocket.Conn
	mu      sync.Mutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]*websocket.Conn),
	}
}

// Register a client connection
func (cm *ClientManager) Register(clientID string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients[clientID] = conn
}

// Unregister a client connection
func (cm *ClientManager) Unregister(clientID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if conn, ok := cm.clients[clientID]; ok {
		conn.Close()
		delete(cm.clients, clientID)
	}
}

// Notify a client
func (cm *ClientManager) Notify(clientID, message string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if conn, ok := cm.clients[clientID]; ok {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Printf("Error notifying client %s: %v", clientID, err)
			conn.Close()
			delete(cm.clients, clientID)
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	clientManager := NewClientManager()
	router := gin.Default()

	// WebSocket endpoint
	router.GET("/ws/:clientID", func(c *gin.Context) {
		clientID := c.Param("clientID")
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}

		clientManager.Register(clientID, conn)
		defer clientManager.Unregister(clientID)

		// Keep the connection open
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Client disconnected: %s, error: %v", clientID, err)
				break
			}
		}
	})

	// File upload endpoint
	router.POST("/upload/:clientID", func(c *gin.Context) {
		clientID := c.Param("clientID")

		// Receive the file
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
			return
		}
		log.Printf("File %s uploaded by client %s", file.Filename, clientID)

		// Save file temporarily (optional)
		dst := fmt.Sprintf("./%s", file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Respond immediately
		c.JSON(http.StatusOK, gin.H{"message": "File received successfully"})

		// Process the file asynchronously
		go func() {
			// Simulate file processing
			time.Sleep(5 * time.Second)

			// Notify the client when processing is complete
			clientManager.Notify(clientID, fmt.Sprintf("File %s processed successfully", file.Filename))
		}()
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
