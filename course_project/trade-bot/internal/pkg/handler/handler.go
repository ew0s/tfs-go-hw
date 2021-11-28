package handler

import (
	"trade-bot/internal/pkg/service"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services   *service.Service
	validate   *validator.Validate
	wsUpgrader *websocket.Upgrader
}

func NewHandler(services *service.Service, validate *validator.Validate, wsUpgrader *websocket.Upgrader) *Handler {
	return &Handler{services: services, validate: validate, wsUpgrader: wsUpgrader}
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

	orderManager := router.Group("/orderManager", h.userIdentity)
	{
		orderManager.POST("send-order", h.sendOrder)
		orderManager.GET("ws/start-trade", h.startTrade)
	}

	return router
}
