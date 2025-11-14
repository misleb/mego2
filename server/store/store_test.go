package store

import (
	"context"
	"os"
	"testing"

	"github.com/misleb/mego2/shared/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDB *DB

func TestMain(m *testing.M) {
	testDBURL, ok := os.LookupEnv("TEST_DATABASE_URL")
	if !ok {
		panic("TEST_DATABASE_URL is not set")
	}

	oldDBURL := os.Getenv("DATABASE_URL")
	os.Setenv("DATABASE_URL", testDBURL)
	defer os.Setenv("DATABASE_URL", oldDBURL)

	var err error
	testDB, err = InitDB()
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	code := m.Run()

	os.Exit(code)
}

func setupTest(t *testing.T) func() {
	testDB.CleanupData()
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, crypt($3, gen_salt('bf')))
		RETURNING id
	`
	var userID int
	err := testDB.db.QueryRow(query, "testuser", "test@example.com", "testpass").Scan(&userID)
	require.NoError(t, err)

	var tokenID int
	err = testDB.db.QueryRow("INSERT INTO tokens (token, user_id) VALUES ($1, $2) RETURNING id", "testtoken", userID).Scan(&tokenID)
	require.NoError(t, err)

	return func() {
		testDB.CleanupData()
	}
}

func TestGetUserByEmailAndPassword(t *testing.T) {
	setupTest(t)
	user, err := testDB.GetUserByEmailAndPassword(context.Background(), "test@example.com", "testpass")

	assert.NoError(t, err)
	assert.NotEmpty(t, user)
}

func TestGetTokenByNameAndPassword_InvalidPassword(t *testing.T) {
	setupTest(t)
	user, err := testDB.GetUserByEmailAndPassword(context.Background(), "test@example.com", "wrongpass")

	assert.Error(t, err)
	assert.Empty(t, user)
}

func TestFindOrCreateUserByEmail_ExistingUser(t *testing.T) {
	setupTest(t)
	user := &types.User{Email: "test@example.com"}
	err := testDB.FindOrCreateUserByEmail(context.Background(), user)
	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.False(t, user.SetPassword)
}

func TestFindOrCreateUserByEmail_NewUser(t *testing.T) {
	setupTest(t)
	user := &types.User{Email: "notfound@example.com", Name: "testuser2"}
	err := testDB.FindOrCreateUserByEmail(context.Background(), user)
	require.NoError(t, err)
	assert.Equal(t, "testuser2", user.Name)
	assert.Equal(t, "notfound@example.com", user.Email)
	assert.True(t, user.SetPassword)

	// then test that the password is not set to blank (that would be bad)
	_, err = testDB.GetUserByEmailAndPassword(context.Background(), "notfound@example.com", "")
	assert.Error(t, err)
}

func TestGetUserByToken_InvalidToken(t *testing.T) {
	setupTest(t)
	user, err := testDB.GetUserByToken(context.Background(), "invalid-token")

	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestGetUserByToken_ValidToken(t *testing.T) {
	setupTest(t)
	user, err := testDB.GetUserByToken(context.Background(), "testtoken")

	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestUpdateUser(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	ctx := context.Background()

	user, err := testDB.GetUserByEmailAndPassword(ctx, "test@example.com", "testpass")
	require.NoError(t, err)
	require.NotNil(t, user)

	user.Password = "newpass"
	err = testDB.UpdateUser(ctx, user, []types.UserColumn{types.UserColPassword})
	require.NoError(t, err)

	// Old password should no longer work
	_, err = testDB.GetUserByEmailAndPassword(ctx, "test@example.com", "testpass")
	assert.Error(t, err)

	// New password should authenticate successfully
	updatedUser, err := testDB.GetUserByEmailAndPassword(ctx, "test@example.com", "newpass")
	require.NoError(t, err)
	assert.Equal(t, user.ID, updatedUser.ID)
}
