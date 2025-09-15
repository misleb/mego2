package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/echo", echoHandler)
	router.NoRoute(gin.WrapH(http.FileServer(http.Dir("./web"))))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	router.Run(":" + port)
}

func echoHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello, World!"})
}
