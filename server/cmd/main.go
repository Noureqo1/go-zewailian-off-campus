package main

import (
	"log"
	"server/db"
	"server/internal/message"
	"server/internal/oauth"
	"server/internal/user"
	"server/internal/ws"
	"server/router"

	"github.com/joho/godotenv"
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

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	messageRepo := message.NewPostgresRepository(dbConn.GetDB())
	var messageCache message.Cache
	if redisClient != nil {
		messageCache = message.NewRedisCache("localhost:6379", "", 0)
	}
	messageSvc := message.NewService(messageRepo, messageCache)
	messageHandler := message.NewHandler(messageSvc)

	hub := ws.NewHub()

	messageAdapter := ws.NewMessageServiceAdapter(messageSvc)
	wsHandler := ws.NewHandler(hub, messageAdapter)

	go hub.Run()

	if wsHandler != nil && messageSvc != nil {
	}

	router.InitRouter(userHandler, wsHandler, messageHandler)
	log.Println("Starting server on :8080")
	if err := router.Start(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

}
