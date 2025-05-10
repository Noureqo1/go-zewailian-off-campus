package router

import (
	"server/internal/auth"
	"server/internal/message"
	"server/internal/user"
	"server/internal/ws"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, wsHandler *ws.Handler, messageHandler *message.Handler, authHandler *auth.Handler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	// Auth routes
	authGroup := r.Group("/auth")
	{
		// Regular auth
		authGroup.POST("/signup", authHandler.Signup)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.GET("/me", authHandler.GetMe)

		// Google auth
		authGroup.GET("/google/login", authHandler.GoogleLogin)
		authGroup.GET("/google/callback", authHandler.GoogleCallback)
	}

	// User routes
	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)

	// WebSocket routes
	r.POST("/ws/createRoom", wsHandler.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/getRooms", wsHandler.GetRooms)
	r.GET("/ws/getClients/:roomId", wsHandler.GetClients)

	// Message API routes
	messageRoutes := r.Group("/api/messages")
	{
		messageRoutes.GET("/room/:roomId", messageHandler.GetMessages)
		messageRoutes.GET("/:messageId", messageHandler.GetMessage)
	}
	
	// Room API routes
	roomRoutes := r.Group("/api/rooms")
	{
		roomRoutes.POST("/", messageHandler.CreateRoom)
		roomRoutes.GET("/", messageHandler.GetRooms)
		roomRoutes.GET("/:roomId", messageHandler.GetRoom)
	}
}

func Start(addr string) error {
	return r.Run(addr)
}
