package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/misleb/mego2/server/endpoint"
	"github.com/misleb/mego2/server/store"
	"github.com/misleb/mego2/shared"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	if err := store.InitDB(); err != nil {
		log.Fatal("Failed to initialize db:", err)
	}
	defer store.CloseDB()

	router := gin.Default()

	endpoint.RegisterEndpoint(router, shared.IncEndpoint, incHandler)
	endpoint.RegisterEndpoint(router, shared.LoginEndpoint, loginHandler)
	router.NoRoute(gin.WrapH(http.FileServer(http.Dir("./web"))))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	if os.Getenv("ENV") != "production" {
		go runWebSocketServer()
	}

	log.Printf("Database URL: %s", os.Getenv("DATABASE_URL"))

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

func incHandler(c *gin.Context, param shared.IntRequest) {
	c.JSON(200, shared.IntResponse{Result: param.Value + 1})
}

func loginHandler(c *gin.Context, param shared.LoginRequest) {
	token, err := store.GetTokenByUser(param.Username, param.Password)
	if err != nil {
		log.Println("login error:", err)
		c.JSON(401, shared.LoginResponse{Error: "Invalid username or password"})
		return
	}
	c.JSON(200, shared.LoginResponse{Token: token})
}
