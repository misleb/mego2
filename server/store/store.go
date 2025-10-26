package store

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/misleb/mego2/shared"
)

type userHash map[string]*shared.User
type tokenHash map[string]*shared.User

var (
	allUsers = userHash{
		"admin":  &shared.User{Name: "admin", Password: "admin"},
		"misleb": &shared.User{Name: "misleb", Password: "kiavfd123"},
	}
	authenticatedUsers = make(tokenHash)
)

func GetUserByToken(token string) *shared.User {
	return authenticatedUsers[token]
}

func GetTokenByUser(name string, pass string) (string, error) {
	if user, exists := allUsers[name]; exists {
		if user.Password == pass {
			token := uuid.New().String()
			authenticatedUsers[token] = user
			return token, nil
		}
		return "", fmt.Errorf("invalid password")
	}
	return "", fmt.Errorf("not found")
}
