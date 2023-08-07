package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mrkouhadi/go-repository-pattern-postgresql/app"
	"github.com/mrkouhadi/go-repository-pattern-postgresql/car"
)

func main() {
	userName := goDotEnvVariable("USER_NAME")
	dbName := goDotEnvVariable("DB_NAME")
	log.Println(userName, dbName)
	dbpool, err := pgxpool.New(context.Background(), "postgres://"+userName+":@localhost:5432/"+dbName)
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

// return the value of the key
func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
