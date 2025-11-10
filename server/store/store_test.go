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
	if err := testutil.SetupTestDB(); err != nil {
		panic(err)
	}

	// Initialize your store with the test DB
	if err := InitDB("../migrations"); err != nil {
		panic(err)
	}

	code := m.Run()

	CloseDB()
	testutil.CleanupTestDB()

	os.Exit(code)
}

func setupTest(t *testing.T) func() {
	testutil.CleanupData(t, db)
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, crypt($3, gen_salt('bf')))
		RETURNING id
	`
	var userID int
	err := db.QueryRow(query, "testuser", "test@example.com", "testpass").Scan(&userID)
	require.NoError(t, err)

	var tokenID int
	err = db.QueryRow("INSERT INTO tokens (token, user_id) VALUES ($1, $2) RETURNING id", "testtoken", userID).Scan(&tokenID)
	require.NoError(t, err)

	return func() {
		testutil.CleanupData(t, db)
	}
}

func TestGetUserByEmailAndPassword(t *testing.T) {
	setupTest(t)
	user, err := GetUserByEmailAndPassword(context.Background(), "test@example.com", "testpass")

	assert.NoError(t, err)
	assert.NotEmpty(t, user)
}

func TestGetTokenByNameAndPassword_InvalidPassword(t *testing.T) {
	setupTest(t)
	user, err := GetUserByEmailAndPassword(context.Background(), "test@example.com", "wrongpass")

	assert.Error(t, err)
	assert.Empty(t, user)
}

func TestFindOrCreateUserByEmail(t *testing.T) {
	setupTest(t)
	// First create a token
	user, err := GetUserByEmailAndPassword(context.Background(), "test@example.com", "testpass")
	require.NoError(t, err)

	// Then retrieve user by token
	user, err = GetUserByToken(context.Background(), user.CurrentToken)
	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestGetUserByToken_InvalidToken(t *testing.T) {
	setupTest(t)
	user, err := GetUserByToken(context.Background(), "invalid-token")

	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestGetUserByToken_ValidToken(t *testing.T) {
	setupTest(t)
	user, err := GetUserByToken(context.Background(), "testtoken")

	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
}
