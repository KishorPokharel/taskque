package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	dbdsn := os.Getenv("DB_DSN")
	db, err := connectDB(dbdsn)
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		db:     db,
		logger: log.Default(),
	}
	if err := app.run(); err != nil {
		log.Fatal(err)
	}
}

func connectDB(dbdsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbdsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, err
}
