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

type DB struct {
	db *sqlx.DB
}

var (
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
)

func (d *DB) Close() error {
	if d.db == nil {
		return nil
	}
	return d.db.Close()
}

func InitDB() (*DB, error) {
	url, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &DB{db: db}, nil
}

// CleanupData removes test data from all tables
func (d *DB) CleanupData() {
	// Clean in order respecting foreign keys
	d.db.Exec("DELETE FROM tokens")
	d.db.Exec("DELETE FROM users")
}

func (d *DB) GetUserByToken(ctx context.Context, token string) (*types.User, error) {
	var user types.User
	tokenModel := &types.Token{Token: token}

	if err := orm.Find(&user).Join(tokenModel).Where("tokens.token = :token").Using(tokenModel).Query(ctx, d.db); err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *DB) GetUserByEmailAndPassword(ctx context.Context, email string, pass string) (*types.User, error) {
	user := &types.User{Email: email, Password: pass}

	// NOTE: Password is not hashed here because it is already hashed in the database and the WHERE will be filtered by the User BeforeFind callback
	// We are absolutely NOT storing the password in the database in plain text. This was an unnecessary abstraction, but fun to play with.
	// :password will be substituted for a PG crypt function call.
	scope := orm.Find(user).Where("email = :email AND password = :password")
	if err := scope.Query(ctx, d.db); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := d.setUserToken(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Used after successful Google authentication to find or create a user with a fresh token
// Caller should prepopulate the user's email and name (from Google)
func (d *DB) FindOrCreateUserByEmail(ctx context.Context, user *types.User) error {
	scope := orm.Find(user).Where("email = :email")
	if err := scope.Query(ctx, d.db); err != nil {
		if err == sql.ErrNoRows {
			user.Password = uuid.New().String() // Random password for new users. This is not a secure password, but it is a temporary password.
			user.IsNewExternal = true
			err := orm.Insert(user).Query(ctx, d.db)
			if err != nil {
				return fmt.Errorf("could not create user: %w", err)
			}
			return d.setUserToken(ctx, user)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}
	return d.setUserToken(ctx, user)
}

func (d *DB) setUserToken(ctx context.Context, user *types.User) error {
	user.CurrentToken = uuid.New().String()
	token := &types.Token{
		Token:  user.CurrentToken,
		UserID: user.ID,
	}
	if err := orm.Insert(token).Query(ctx, d.db); err != nil {
		return fmt.Errorf("could not create token: %w", err)
	}
	return nil
}

func (d *DB) UpdateUser(ctx context.Context, user *types.User, columns []types.UserColumn) error {
	return orm.Update(user).Set(d.convertUserColumnsToStrings(columns)).Query(ctx, d.db)
}

// This is a helper function to convert the UserColumn enum to a string slice for the UpdateModel
// This isn't strictly necessary, but it helps to enforce that the columns are valid for the User model at compile time.
func (d *DB) convertUserColumnsToStrings(columns []types.UserColumn) []string {
	columnStrings := make([]string, len(columns))
	for i, column := range columns {
		columnStrings[i] = string(column)
	}
	return columnStrings
}
