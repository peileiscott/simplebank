package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const dsn = "postgresql://postgres:password@localhost:5432/simplebank?sslmode=disable"

var testQueries *Queries

func TestMain(m *testing.M) {
	pool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(pool)
	os.Exit(m.Run())
}
