package main

import (
	"fmt"
	controllers "go-image-bg-remover/controllers"
	middlewares "go-image-bg-remover/middlewares"
	"go-image-bg-remover/socket"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Allow all origins for simplicity. Adjust as needed for security.
		return true
	},
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	clientManager := socket.NewClientManager()

	router := gin.Default()

	// WebSocket endpoint
	router.GET("/ws/:clientID", func(c *gin.Context) {
		clientID := c.Param("clientID")
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Printf("Failed to upgrade connection: %v", err)
			return
		}

		clientManager.Register(clientID, conn)
		defer clientManager.Unregister(clientID)

		// Keep the connection open
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Client disconnected: %s, error: %v", clientID, err)
				break
			}
		}
	})

	router.POST("/upload", controllers.UploadImage)
	router.GET("/download", middlewares.ValidateSignedURLMiddleware(), controllers.DownloadImage)

	router.Run(":" + port)
}
