package main

import (
	"log"
	"server/db"
	"server/internal/user"
	"server/router"
)

func main() { 
	log.Println("Starting the application...")
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %v", err)
	}
	log.Println("Database connection established")

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	router.InitRouter(userHandler)
	log.Println("Starting server on :8080")
	if err := router.Start(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
