package handler

import (
	"chat/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	api := router.Group("/api", h.userIdentity)
	{
		messages := api.Group("/messages")
		{
			messages.POST("/", h.createGlobalMessage)
			messages.GET("/", h.getGlobalMessages)
		}
		users := api.Group("/users")
		{
			users.POST("/:id/messages", h.sendMessageToUserByID)
			users.GET("/messages", h.getUserMessages)
		}
	}

	return router
}
