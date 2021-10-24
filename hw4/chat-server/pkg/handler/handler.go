package handler

import (
	"chat/pkg/service"

	_ "chat/docs" // docs for swagger

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
