package store

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/misleb/mego2/shared/types"
)

var (
	db                 *sql.DB
	GoogleClientSecret = func() string {
		secret, ok := os.LookupEnv("GOOGLE_CLIENT_SECRET")
		if !ok {
			panic("GOOGLE_CLIENT_SECRET is not set")
		}
		return secret
	}()
	GoogleClientID = func() string {
		id, ok := os.LookupEnv("GOOGLE_CLIENT_ID")
		if !ok {
			panic("GOOGLE_CLIENT_ID is not set")
		}
		return id
	}()
)

func InitDB() error {
	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return fmt.Errorf("no DATABASE_URL set")
	}

	var err error
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return db.Ping()
}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func GetUserByToken(token string) *types.User {
	query := `SELECT users.name, users.email FROM users LEFT JOIN tokens ON users.id = tokens.user_id WHERE tokens.token = $1`

	var user types.User

	err := db.QueryRow(query, token).Scan(&user.Name, &user.Email)
	if err != nil {
		return nil
	}
	return &user
}

func GetTokenByUser(name string, pass string) (string, error) {
	user, err := fetchUserAndToken(name, pass)
	if err == nil {
		return user.Token, nil
	}
	return "", err
}

func fetchUserAndToken(name string, pass string) (*types.User, error) {
	uQuery := `SELECT id, email, name FROM users WHERE name = $1 AND crypt($2, password) = password`
	tQuery := `INSERT INTO tokens (token, user_id) VALUES ($1, $2) RETURNING id`

	var user types.User
	var id int32

	if err := db.QueryRow(uQuery, name, pass).Scan(&id, &user.Email, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Token = uuid.New().String()

	if err := db.QueryRow(tQuery, user.Token, id).Scan(&id); err != nil {
		return nil, fmt.Errorf("could not create token: %w", err)
	}
	return &user, nil
}
