package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/misleb/mego2/server/endpoint"
	"github.com/misleb/mego2/shared"
)

func main() {
	router := gin.Default()

	endpoint.RegisterEndpoint(router, shared.IncEndpoint, incHandler)
	endpoint.RegisterEndpoint(router, shared.LoginEndpoint, loginHandler)
	router.NoRoute(gin.WrapH(http.FileServer(http.Dir("./web"))))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	router.Run(":" + port)
}

func incHandler(c *gin.Context, param shared.IntRequest) {
	c.JSON(200, shared.IntResponse{Result: param.Value + 1})
}

func loginHandler(c *gin.Context, param shared.LoginRequest) {
	if param.Username != "admin" || param.Password != "admin" {
		c.JSON(401, shared.LoginResponse{Error: "Invalid username or password"})
		return
	}
	token := uuid.New().String()
	c.JSON(200, shared.LoginResponse{Token: token})
}
