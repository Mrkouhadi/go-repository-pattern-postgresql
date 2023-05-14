package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrkouhadi/go-repository-pattern-postgresql/app"
	"github.com/mrkouhadi/go-repository-pattern-postgresql/car"
)

func main() {

	dbpool, err := pgxpool.New(context.Background(), "postgres://kouhadi:@localhost:5432/go-pgx-repository")
	// url example: "postgres://username:password@localhost:5432/database_name"
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	carRepo := car.NewPgxRepository(dbpool)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.RunRepository(ctx, carRepo)
}
