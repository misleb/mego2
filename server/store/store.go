package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/misleb/mego2/shared/orm"
	"github.com/misleb/mego2/shared/types"
)

var (
	db                 *sqlx.DB
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
	BaseURI = func() string {
		id, ok := os.LookupEnv("BASE_URI")
		if !ok {
			panic("BASE_URI is not set")
		}
		return id
	}
	databaseURL = func() string {
		url, ok := os.LookupEnv("DATABASE_URL")
		if !ok {
			panic("DATABASE_URL is not set")
		}
		return url
	}()
)

func InitDB() error {
	var err error
	db, err = sqlx.Open("postgres", databaseURL)
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

func GetUserByToken(ctx context.Context, token string) (*types.User, error) {
	var user types.User
	tokenModel := &types.Token{Token: token}

	if err := orm.Find(&user).Join(tokenModel).Where("tokens.token = :token").Using(tokenModel).Query(ctx, db); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetTokenByNameAndPassword(ctx context.Context, name string, pass string) (string, error) {
	user, err := fetchUserAndToken(ctx, name, pass)
	if err == nil {
		return user.CurrentToken, nil
	}
	return "", err
}

func fetchUserAndToken(ctx context.Context, name string, pass string) (*types.User, error) {
	user := &types.User{Name: name, Password: pass}

	scope := orm.Find(user).Where("name = :name AND crypt(:password, password) = password") // TODO: use sqlx to bind the password with a func
	if err := scope.Query(ctx, db); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := setUserToken(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Used after successful Google authentication to find or create a user with a fresh token
// Caller should prepopulate the user's email and name (from Google)
func FindOrCreateUserByEmail(ctx context.Context, user *types.User) error {
	scope := orm.Find(user).Where("email = :email")
	if err := scope.Query(ctx, db); err != nil {
		if err == sql.ErrNoRows {
			user.Password = "test" // TODO: remove this
			err := orm.Insert(user).Query(ctx, db)
			if err != nil {
				return fmt.Errorf("could not create user: %w", err)
			}
			return setUserToken(ctx, user)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}
	return setUserToken(ctx, user)
}

func setUserToken(ctx context.Context, user *types.User) error {
	user.CurrentToken = uuid.New().String()
	token := &types.Token{
		Token:  user.CurrentToken,
		UserID: user.ID,
	}
	if err := orm.Insert(token).Query(ctx, db); err != nil {
		return fmt.Errorf("could not create token: %w", err)
	}
	return nil
}
