package main

import (
	"log"
	"server/db"
	"server/internal/message"
	"server/internal/oauth"
	"server/internal/user"
	"server/internal/ws"
	"server/router"
	"time"

	"github.com/joho/godotenv"
	"github.com/sony/gobreaker"
	"server/internal/auth"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or failed to load .env")
	}

	log.Println("Starting the application...")
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %v", err)
	}
	log.Println("Database connection established")

	redisClient, err := db.NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v. Proceeding without caching.", err)
	} else {
		log.Println("Redis connection established")
		defer redisClient.Close()
	}

	oauth.InitGoogleOAuth()

	// Initialize auth components
	authRepo := auth.NewRepository(dbConn.GetDB())
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	messageRepo := message.NewPostgresRepository(dbConn.GetDB())
	var messageCache message.Cache
	if redisClient != nil {
		messageCache = message.NewRedisCache("localhost:6379", "", 0)
	}
	// Create base message service
	baseSvc := message.NewService(messageRepo, messageCache)

	// Configure circuit breaker
	cbConfig := message.CircuitBreakerConfig{
		Name:        "chat-service",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit breaker %s state changed from %s to %s", name, from, to)
		},
	}

	// Configure retry mechanism
	retryConfig := message.RetryConfig{
		MaxElapsedTime:   1 * time.Minute,
		MaxInterval:      5 * time.Second,
		InitialInterval: 100 * time.Millisecond,
	}

	// Create resilient service wrapper
	messageSvc := message.NewResilientService(baseSvc, cbConfig, retryConfig)
	messageHandler := message.NewHandler(messageSvc)

	hub := ws.NewHub()

	messageAdapter := ws.NewMessageServiceAdapter(messageSvc)
	wsHandler := ws.NewHandler(hub, messageAdapter)

	go hub.Run()

	// Update router initialization to include auth handler
	router.InitRouter(userHandler, wsHandler, messageHandler, authHandler)
	log.Println("Starting server on :8080")
	if err := router.Start(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
