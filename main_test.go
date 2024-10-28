package pgxfactory

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

var (
	pool *pgxpool.Pool
	ctx  = context.Background()
)

func TestMain(m *testing.M) {
	p, err := pgxpool.New(ctx, "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer p.Close()

	pool = p

	code := m.Run()
	os.Exit(code)
}
