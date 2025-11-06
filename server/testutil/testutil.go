// server/testutil/testutil.go
package testutil

import (
	"database/sql"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

var (
	TestDB        *sqlx.DB
	TestDBURL     string
	dbInitialized bool
)

// SetupTestDB initializes the test database and runs migrations
// Call this from TestMain in each package that needs it
func SetupTestDB(migrationsPath string) error {
	if dbInitialized {
		return nil // Already set up
	}

	TestDBURL = os.Getenv("TEST_DATABASE_URL")
	if TestDBURL == "" {
		return nil // Tests will be skipped
	}

	os.Setenv("DATABASE_URL", TestDBURL)

	var err error
	TestDB, err = sqlx.Open("postgres", TestDBURL)
	if err != nil {
		return err
	}

	if err := runMigrations(TestDBURL, migrationsPath); err != nil {
		return err
	}

	dbInitialized = true
	return nil
}

func CleanupTestDB() {
	if TestDB != nil {
		TestDB.Close()
	}
}

func runMigrations(dbURL, migrationsPath string) error {
	migDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer migDB.Close()

	driver, err := postgres.WithInstance(migDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	m.Down()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// CleanupData removes test data from all tables
func CleanupData(t *testing.T, db *sqlx.DB) {
	// Clean in order respecting foreign keys
	db.Exec("DELETE FROM tokens")
	db.Exec("DELETE FROM users")
}
