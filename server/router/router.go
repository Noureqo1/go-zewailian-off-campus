package router

import (
	"server/internal/oauth"
	"server/internal/user"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler) {
	r = gin.Default()

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)

	// Google OAuth2 endpoints
	r.GET("/auth/google/login", gin.WrapF(oauth.GoogleLoginHandler))
	r.GET("/auth/google/callback", gin.WrapF(oauth.GoogleCallbackHandler))
}

func Start(addr string) error {
	return r.Run(addr)
}
