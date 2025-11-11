package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/misleb/mego2/server/endpoint"
	"github.com/misleb/mego2/server/store"
	"github.com/misleb/mego2/shared/types"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	if err := store.InitDB("./server/migrations"); err != nil {
		log.Fatal("Failed to initialize db:", err)
	}
	defer store.CloseDB()

	router := gin.Default()

	endpoint.RegisterEndpoint(router, types.IncEndpoint, incHandler)
	endpoint.RegisterEndpoint(router, types.LoginEndpoint, loginHandler)
	endpoint.RegisterEndpoint(router, types.GoogleAuthEndpoint, googleAuthHandler)
	endpoint.RegisterEndpoint(router, types.UpdateSelfEndpoint, updateSelfHandler)
	router.NoRoute(gin.WrapH(http.FileServer(http.Dir("./web"))))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	if os.Getenv("ENV") != "production" {
		go runWebSocketServer()
	}

	router.Run(":" + port)
}

func runWebSocketServer() {
	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to WebSocket:", err)
			return
		}
		defer ws.Close()

		// Handle WebSocket messages
		for {
			messageType, msg, err := ws.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}

			// Echo the message back
			err = ws.WriteMessage(messageType, msg)
			if err != nil {
				log.Println("Error writing message:", err)
				break
			}
		}
	})

	router.Run(":38919") // Different port for WebSocket server
}
