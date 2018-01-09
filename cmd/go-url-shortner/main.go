package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	// Postgres connector
	_ "github.com/lib/pq"
	urlshortner "github.com/mauricioklein/go-url-shortner"
	"github.com/mauricioklein/go-url-shortner/store"
)

var (
	errNoDbHost     = errors.New("no postgres host provided")
	errNoDbUser     = errors.New("no postgres user provided")
	errNoDbDatabase = errors.New("no postgres database provided")
)

func main() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		panic(errNoDbHost)
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		panic(errNoDbUser)
	}

	database := os.Getenv("DB_DATABASE")
	if database == "" {
		panic(errNoDbDatabase)
	}

	dbinfo := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", host, user, database)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	linkStore := store.NewLinkStore(db)

	muxRouter := urlshortner.NewRouter(linkStore)
	http.ListenAndServe(":8080", muxRouter)
}
