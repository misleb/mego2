package main

import (
	"embed"
	"io/fs"
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

//go:embed web/*
var webFS embed.FS

func main() {
	db, err := store.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize db:", err)
	}
	defer db.Close()

	webContent, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatal("Failed to get web content:", err)
	}

	router := gin.Default()

	endpoint.RegisterEndpoint(router, types.IncEndpoint, incHandler, db)
	endpoint.RegisterEndpoint(router, types.LoginEndpoint, loginHandler, db)
	endpoint.RegisterEndpoint(router, types.GoogleAuthEndpoint, googleAuthHandler, db)
	endpoint.RegisterEndpoint(router, types.UpdateSelfEndpoint, updateSelfHandler, db)
	router.NoRoute(gin.WrapH(http.FileServer(http.FS(webContent))))

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
