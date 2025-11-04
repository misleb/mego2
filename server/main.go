package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/misleb/mego2/server/endpoint"
	"github.com/misleb/mego2/server/store"
	"github.com/misleb/mego2/shared/api_client"
	"github.com/misleb/mego2/shared/types"
	"google.golang.org/api/idtoken"
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

	endpoint.RegisterEndpoint(router, types.IncEndpoint, incHandler)
	endpoint.RegisterEndpoint(router, types.LoginEndpoint, loginHandler)
	endpoint.RegisterEndpoint(router, types.GoogleAuthEndpoint, googleAuthHandler)
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

func incHandler(c *gin.Context, param types.IntRequest) {
	c.JSON(200, types.IntResponse{Result: param.Value + 1})
}

func loginHandler(c *gin.Context, param types.LoginRequest) {
	token, err := store.GetTokenByUser(param.Username, param.Password)
	if err != nil {
		log.Println("login error:", err)
		c.JSON(401, types.LoginResponse{Error: "Invalid username or password"})
		return
	}
	c.JSON(200, types.LoginResponse{Token: token})
}

func googleAuthHandler(c *gin.Context, param types.GoogleAuthRequest) {
	remote := api_client.GetInstance()
	request := types.GoogleTokenExchangeRequest{
		Code:         param.Code,
		ClientID:     store.GoogleClientID,
		ClientSecret: store.GoogleClientSecret,
		RedirectURI:  store.BaseURI() + "/google-callback.html",
		GrantType:    "authorization_code",
	}

	tokenExchangeResponse, err := api_client.CallEndpointTyped[types.GoogleTokenExchangeResponse](
		remote, types.GoogleTokenExchangeEndpoint, request, api_client.NoOpRequestAugment,
	)
	if err != nil {
		log.Println("google error:", err)
		c.JSON(401, types.GoogleAuthResponse{Error: "Couldn't validate with Google"})
		return
	}

	payload, err := idtoken.Validate(context.Background(), tokenExchangeResponse.IDToken, store.GoogleClientID)
	if err != nil {
		log.Println("ID token validation error:", err)
		c.JSON(401, types.GoogleAuthResponse{Error: "Invalid ID token"})
		return
	}

	// Extract user info from the token payload
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	// TODO: Create or update user in database, generate app token

	c.JSON(200, types.GoogleAuthResponse{
		Token: "test", // TODO: Generate actual app token
		Email: email,
		Name:  name,
	})
}
