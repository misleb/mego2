package store

import (
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
	query := `
		INSERT INTO users (name, email, password) 
		VALUES ($1, $2, crypt($3, gen_salt('bf')))
		ON CONFLICT DO NOTHING
	`
	_, err := db.Exec(query, "testuser", "test@example.com", "testpass")
	require.NoError(t, err)

	return func() {
		testutil.CleanupData(t, db)
	}
}

func TestGetTokenByUser(t *testing.T) {
	setupTest(t)
	token, err := GetTokenByUser("testuser", "testpass")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGetTokenByUser_InvalidPassword(t *testing.T) {
	setupTest(t)
	token, err := GetTokenByUser("testuser", "wrongpass")

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestGetUserByToken(t *testing.T) {
	setupTest(t)
	// First create a token
	token, err := GetTokenByUser("testuser", "testpass")
	require.NoError(t, err)

	// Then retrieve user by token
	user := GetUserByToken(token)

	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestGetUserByToken_InvalidToken(t *testing.T) {
	setupTest(t)
	user := GetUserByToken("invalid-token")

	// Current implementation returns empty user, not nil
	assert.Equal(t, "", user.Name)
	assert.Equal(t, "", user.Email)
}
