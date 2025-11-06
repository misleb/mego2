package store

import (
	"context"
	"os"
	"testing"

	"github.com/misleb/mego2/server/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if err := testutil.SetupTestDB("../migrations"); err != nil {
		panic(err)
	}

	// Initialize your store with the test DB
	if err := InitDB(); err != nil {
		panic(err)
	}

	code := m.Run()

	CloseDB()
	testutil.CleanupTestDB()

	os.Exit(code)
}

func setupTest(t *testing.T) func() {
	_, err := db.ExecContext(context.Background(), "DELETE FROM tokens")
	require.NoError(t, err)
	_, err = db.ExecContext(context.Background(), "DELETE FROM users")
	require.NoError(t, err)
	query := `
		INSERT INTO users (name, email, password) 
		VALUES ($1, $2, crypt($3, gen_salt('bf')))
		RETURNING id
	`
	var userID int
	err = db.QueryRow(query, "testuser", "test@example.com", "testpass").Scan(&userID)
	require.NoError(t, err)

	_, err = db.ExecContext(context.Background(), "INSERT INTO tokens (token, user_id) VALUES ($1, $2)", "testtoken", userID)
	require.NoError(t, err)

	return func() {
		testutil.CleanupData(t, db)
	}
}

func TestGetTokenByNameAndPassword(t *testing.T) {
	setupTest(t)
	token, err := GetTokenByNameAndPassword(context.Background(), "testuser", "testpass")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGetTokenByNameAndPassword_InvalidPassword(t *testing.T) {
	setupTest(t)
	token, err := GetTokenByNameAndPassword(context.Background(), "testuser", "wrongpass")

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestFindOrCreateUserByEmail(t *testing.T) {
	setupTest(t)
	// First create a token
	token, err := GetTokenByNameAndPassword(context.Background(), "testuser", "testpass")
	require.NoError(t, err)

	// Then retrieve user by token
	user := GetUserByToken(context.Background(), token)

	require.NotNil(t, user)
	assert.Equal(t, "testuser", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestGetUserByToken_InvalidToken(t *testing.T) {
	setupTest(t)
	user := GetUserByToken(context.Background(), "invalid-token")

	assert.Nil(t, user)
}

func TestGetUserByToken_ValidToken(t *testing.T) {
	setupTest(t)
	user := GetUserByToken(context.Background(), "testtoken")

	require.NotNil(t, user)
	assert.Equal(t, "testuser", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
}
