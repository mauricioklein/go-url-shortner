package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	// Postgres connector
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	urlshortner "github.com/mauricioklein/go-url-shortner"
	"github.com/mauricioklein/go-url-shortner/store"
)

var (
	// Postgres errors
	errNoDbHost     = errors.New("no postgres host provided")
	errNoDbUser     = errors.New("no postgres user provided")
	errNoDbDatabase = errors.New("no postgres database provided")

	// Redis errors
	errNoRedisHost = errors.New("no redis host provided")
	errNoRedisPort = errors.New("no redis port provided")
)

func main() {
	pgHost := os.Getenv("DB_HOST")
	if pgHost == "" {
		panic(errNoDbHost)
	}

	pgUser := os.Getenv("DB_USER")
	if pgUser == "" {
		panic(errNoDbUser)
	}

	pgDb := os.Getenv("DB_DATABASE")
	if pgDb == "" {
		panic(errNoDbDatabase)
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		panic(errNoRedisHost)
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		panic(errNoRedisPort)
	}

	// establish connection to Postgres
	dbinfo := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", pgHost, pgUser, pgDb)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	linkStore := store.NewLinkStore(db)

	// establish connection to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	_, err = redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	redisStore := store.NewRedisStore(redisClient)
	muxRouter := urlshortner.NewRouter(linkStore, redisStore)

	fmt.Printf("Now serving on port :8080")
	http.ListenAndServe(":8080", muxRouter)
}
