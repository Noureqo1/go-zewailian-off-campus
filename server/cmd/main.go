package main

import (
	"log"
	"server/db"
	"server/internal/oauth"
	"server/internal/user"
	"server/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or failed to load .env")
	}

	log.Println("Starting the application...")
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %v", err)
	}
	log.Println("Database connection established")

	// Initialize Google OAuth2 config
	oauth.InitGoogleOAuth()

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	router.InitRouter(userHandler)
	log.Println("Starting server on :8080")
	if err := router.Start(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
