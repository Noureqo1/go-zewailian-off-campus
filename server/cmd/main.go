package main

import (
	"log"
	"server/db"
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

	oauth.InitGoogleOAuth()

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)

	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	log.Println("Starting server on :8081")
	if err := router.Start(":8081"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

}
