package store

import (
	"context"
	"os"
	"testing"

	"github.com/misleb/mego2/server/testutil"
	"github.com/misleb/mego2/shared/types"
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

func TestFindOrCreateUserByEmail_ExistingUser(t *testing.T) {
	setupTest(t)
	user := &types.User{Email: "test@example.com"}
	err := FindOrCreateUserByEmail(context.Background(), user)
	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.False(t, user.SetPassword)
}

func TestFindOrCreateUserByEmail_NewUser(t *testing.T) {
	setupTest(t)
	user := &types.User{Email: "notfound@example.com", Name: "testuser2"}
	err := FindOrCreateUserByEmail(context.Background(), user)
	require.NoError(t, err)
	assert.Equal(t, "testuser2", user.Name)
	assert.Equal(t, "notfound@example.com", user.Email)
	assert.True(t, user.SetPassword)

	// then test that the password is not set to blank (that would be bad)
	_, err = GetUserByEmailAndPassword(context.Background(), "notfound@example.com", "")
	assert.Error(t, err)
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

func TestUpdateUser(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	ctx := context.Background()

	user, err := GetUserByEmailAndPassword(ctx, "test@example.com", "testpass")
	require.NoError(t, err)
	require.NotNil(t, user)

	user.Password = "newpass"
	err = UpdateUser(ctx, user, []types.UserColumn{types.UserColPassword})
	require.NoError(t, err)

	// Old password should no longer work
	_, err = GetUserByEmailAndPassword(ctx, "test@example.com", "testpass")
	assert.Error(t, err)

	// New password should authenticate successfully
	updatedUser, err := GetUserByEmailAndPassword(ctx, "test@example.com", "newpass")
	require.NoError(t, err)
	assert.Equal(t, user.ID, updatedUser.ID)
}
