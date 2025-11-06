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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
	token, err := store.GetTokenByNameAndPassword(c.Request.Context(), param.Username, param.Password)
	if err != nil {
		log.Println("login error:", err)
		c.JSON(401, types.LoginResponse{Error: "Invalid username or password"})
		return
	}
	c.JSON(200, types.LoginResponse{Token: token})
}

func googleAuthHandler(c *gin.Context, param types.GoogleAuthRequest) {
	config := &oauth2.Config{
		ClientID:     store.GoogleClientID,
		ClientSecret: store.GoogleClientSecret,
		RedirectURL:  store.BaseURI() + "/google-callback.html",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}

	token, err := config.Exchange(c.Request.Context(), param.Code)
	if err != nil {
		log.Println("token exchange error:", err)
		c.JSON(401, types.GoogleAuthResponse{Error: "Couldn't validate with Google"})
		return
	}

	// Extract ID token from the oauth2.Token
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Println("ID token not found in response")
		c.JSON(401, types.GoogleAuthResponse{Error: "Invalid token response"})
		return
	}

	payload, err := idtoken.Validate(c.Request.Context(), idToken, store.GoogleClientID)
	if err != nil {
		log.Println("ID token validation error:", err)
		c.JSON(401, types.GoogleAuthResponse{Error: "Invalid ID token"})
		return
	}

	// Extract user info from the token payload
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	user := &types.User{Email: email, Name: name}
	err = store.FindOrCreateUserByEmail(c.Request.Context(), user)
	if err != nil {
		log.Println("find or create user error:", err)
		c.JSON(401, types.GoogleAuthResponse{Error: "Couldn't find or create user"})
		return
	}

	c.JSON(200, types.GoogleAuthResponse{
		Token: user.CurrentToken,
		Email: email,
		Name:  name,
	})
}
