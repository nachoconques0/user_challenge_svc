package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/controller/http/user"
)

// InitUserRoutes will set all the endpoints for an user
func InitUserRoutes(
	router *gin.Engine,
	userCtrl *user.Controller,
) {
	userGroup := router.Group("/users")
	userGroup.GET("", userCtrl.Find)
	userGroup.POST("", userCtrl.Create)
	userGroup.PATCH("/:id", userCtrl.Update)
	userGroup.DELETE("/:id", userCtrl.Delete)
}
