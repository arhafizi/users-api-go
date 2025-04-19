package services_test // Or a central test setup package

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var testDB *sql.DB
var testTableNames = []string{"users"}

func TestMain(m *testing.M) {
	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:4321@localhost:5432/usersdb?sslmode=disable"
	}

	var err error
	testDB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = testDB.PingContext(ctx); err != nil {
		testDB.Close()
		log.Fatalf("Failed to ping test database: %v", err)
	}

	log.Println("Test database connection established.")

	exitCode := m.Run()

	log.Println("Closing test database connection.")
	testDB.Close()
	os.Exit(exitCode)
}

func TruncateTables(t *testing.T, db *sql.DB, tables []string) {
	if len(tables) == 0 {
		return
	}

	allowedTables := map[string]bool{"users": true}
	for _, table := range tables {
		if !allowedTables[table] {
			log.Fatalf("Attempted to truncate disallowed table: %s", table)
		}
	}

	// Quote table names to handle potential special characters or reserved words
	quotedTables := make([]string, len(tables))
	for i, table := range tables {
		quotedTables[i] = fmt.Sprintf(`"%s"`, table) // double quotes for standard SQL identifiers
	}
	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(quotedTables, ", "))
	_, err := db.Exec(query)
	require.NoError(t, err, "Failed to truncate tables: %v", tables)
}

func seedUser(t *testing.T, email string, options ...string) int32 {
	ctx := context.Background()
	var id int32

	// Default values
	username := "usr_" + strings.Split(email, "@")[0]
	fullName := "Test User " + strings.Split(email, "@")[0]
	passwordHash := "hashed_password"

	if len(options) > 0 && options[0] != "" {
		username = options[0]
	}
	if len(options) > 1 && options[1] != "" {
		fullName = options[1]
	}
	if len(options) > 2 && options[2] != "" {
		passwordHash = options[2]
	}

	query := `
        INSERT INTO users (username, email, full_name, password_hash) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id`
	err := testDB.QueryRowContext(ctx, query, username, email, fullName, passwordHash).Scan(&id)
	require.NoError(t, err, "Failed to seed user with email: %s", email)
	return id
}
