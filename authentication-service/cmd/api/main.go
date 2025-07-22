package main

import (
	"authentication-service/data"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *pgx.Conn
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// TODO connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Cannot connect to Postgres!")
	}

	// setup config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer db.Close(context.Background())

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *pgx.Conn {
	dsn := os.Getenv("DSN")
	log.Println(dsn)

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
	}
}
