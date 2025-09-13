package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dhruvpatel-10/signee/ca-api/cmd/api"
	"github.com/dhruvpatel-10/signee/ca-api/db"
	"github.com/dhruvpatel-10/signee/ca-api/internal/config"
	"github.com/dhruvpatel-10/signee/ca-api/internal/service/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func initializeApp() {
	config.LoadEnvVariables()
	logger.Init()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fix: Origin should be the frontend domain, not the API path
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func openDatabase() (*sql.DB, *db.Queries, error) {
	dbURL := os.Getenv("GOOSE_DBSTRING")
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, err
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		return nil, nil, err
	}

	return conn, db.New(conn), nil
}

func newRouter(queries *db.Queries) *gin.Engine {
	r := gin.New()

	r.Use(CORSMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api.SetupRoutes(r, queries)
	return r
}

func serveUnixSocket(handler http.Handler, socketPath string) (*http.Server, net.Listener, error) {
	// Remove existing socket if it exists
	if err := os.RemoveAll(socketPath); err != nil {
		return nil, nil, err
	}

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, nil, err
	}

	// Set proper permissions for the socket
	if err := os.Chmod(socketPath, 0666); err != nil {
		return nil, nil, err
	}

	server := &http.Server{Handler: handler}
	log.Println("Listening on unix socket:", socketPath)
	return server, l, nil
}

func main() {
	initializeApp()

	conn, queries, err := openDatabase()
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	router := newRouter(queries)

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	const socketPath = "/tmp/go.sock"

	// Use environment variable to choose socket vs port
	if useSocket := os.Getenv("USE_UNIX_SOCKET"); useSocket == "true" {
		server, listener, err := serveUnixSocket(router, socketPath)
		if err != nil {
			log.Fatalf("unix socket setup error: %v", err)
		}

		// Start server in a goroutine
		go func() {
			if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
				log.Fatalf("unix socket server error: %v", err)
			}
		}()

		// Wait for signal
		<-sigChan
		log.Println("[shutting down] cleaning...")

		// Create context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown server gracefully
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
		}

		// Clean up socket
		os.RemoveAll(socketPath)
		log.Println("[shutting down] see you again~")

	} else {
		port := os.Getenv("GO_PORT")
		if port == "" {
			port = "8080"
		}

		server := &http.Server{
			Addr:    ":" + port,
			Handler: router,
		}

		// Start server in a goroutine
		go func() {
			log.Printf("Server starting on port %s", port)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("http server error: %v", err)
			}
		}()

		// Wait for signal
		<-sigChan
		log.Println("[shutting down] cleaning...")

		// Create context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown server gracefully
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
		}

		log.Println("[shutting down] see you again~")
	}
}
