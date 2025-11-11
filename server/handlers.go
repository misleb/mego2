package main

import (
	"fmt"
	"log"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/gin-gonic/gin"
	"github.com/misleb/mego2/server/store"
	"github.com/misleb/mego2/shared/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func incHandler(c *gin.Context, param types.IntRequest) {
	c.JSON(200, types.IntResponse{Result: param.Value + 1})
}

func loginHandler(c *gin.Context, param types.LoginRequest) {
	user, err := store.GetUserByEmailAndPassword(c.Request.Context(), param.Username, param.Password)
	if err != nil {
		log.Println("login error:", err)
		c.JSON(401, types.LoginResponse{Error: "Invalid username or password"})
		return
	}
	c.JSON(200, types.LoginResponse{User: user})
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
		c.JSON(401, types.LoginResponse{Error: "Couldn't validate with Google"})
		return
	}

	// Extract ID token from the oauth2.Token
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Println("ID token not found in response")
		c.JSON(401, types.LoginResponse{Error: "Invalid token response"})
		return
	}

	payload, err := idtoken.Validate(c.Request.Context(), idToken, store.GoogleClientID)
	if err != nil {
		log.Println("ID token validation error:", err)
		c.JSON(401, types.LoginResponse{Error: "Invalid ID token"})
		return
	}

	// Extract user info from the token payload
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	user := &types.User{Email: email, Name: name}
	err = store.FindOrCreateUserByEmail(c.Request.Context(), user)
	if err != nil {
		log.Println("find or create user error:", err)
		c.JSON(401, types.LoginResponse{Error: "Couldn't find or create user"})
		return
	}
	c.JSON(200, types.LoginResponse{
		User: user,
	})
}

func updateSelfHandler(c *gin.Context, param types.UpdateSelfRequest) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.JSON(401, types.LoginResponse{Error: err.Error()})
		return
	}

	if param.NewPassword != "" {
		if user.IsNewExternal {
			user.Password = param.NewPassword // TODO: check password strength
			user.IsNewExternal = false
		} else {
			if _, err := store.GetUserByEmailAndPassword(c.Request.Context(), user.Email, param.OldPassword); err != nil {
				c.JSON(401, types.LoginResponse{Error: "invalid old password"})
				return
			} else {
				user.Password = param.NewPassword
			}
		}
	}

	// TODO: Support more than password update
	if err := store.UpdateUser(c.Request.Context(), user, []types.UserColumn{types.UserColPassword}); err != nil {
		c.JSON(500, types.LoginResponse{Error: err.Error()})
		return
	} else {
		c.JSON(200, types.LoginResponse{User: user})
	}
}

func getCurrentUser(c *gin.Context) (*types.User, error) {
	if userUntyped, ok := c.Get("currentUser"); !ok {
		return nil, fmt.Errorf("no current user found")
	} else {
		return userUntyped.(*types.User), nil
	}
}
