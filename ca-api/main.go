package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/dhruvpatel-10/signee/ca-api/cmd/api"
	"github.com/dhruvpatel-10/signee/ca-api/db"
	"github.com/dhruvpatel-10/signee/ca-api/internal/config"
	"github.com/dhruvpatel-10/signee/ca-api/internal/service/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func initiliser() {
	config.LoadEnvVariables()
	logger.Init()
}

func main() {
	initiliser()

	db_url := os.Getenv("GOOSE_DBSTRING")
	// open db connection
	conn, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	// wrap *sql.DB with sqlc Queries
	queries := db.New(conn)

	// setup Gin
	r := gin.Default()
	api.SetupRoutes(r, queries)

	// start server
	if err := r.Run(":9000"); err != nil {
		log.Fatal(err)
	}
}
