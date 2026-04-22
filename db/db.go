package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("dsn : ", dsn)
	var err error
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Println("DB connected")
}
