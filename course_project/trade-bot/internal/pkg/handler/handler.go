package handler

import (
	"trade-bot/internal/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	auth := router.Group("/auth")
	{
		auth.POST("sign-in", h.signIn)
		auth.POST("sign-up", h.signUp)
		auth.DELETE("logout", h.userIdentity, h.logout)
	}

	todo := router.Group("/todo", h.userIdentity)
	{
		todo.GET("/", h.getTodo)
	}

	return router
}
