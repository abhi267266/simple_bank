package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	// Setup DB connection pool
	DBurl := "postgresql://root:secret@localhost:5050/simple_bank?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), DBurl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	testDB = pool
	testQueries = New(pool)
	defer testDB.Close() // Defer close of pool after tests complete

	// Run tests
	code := m.Run()

	// Exit with the appropriate exit code
	os.Exit(code)
}
