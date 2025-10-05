package handler

import (
	talk_together_app "github.com/firstproject/talk-together-app/hub"
	"github.com/firstproject/talk-together-app/pkg/service"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "github.com/firstproject/talk-together-app/docs"
)

type Handler struct {
	services *service.Service
	hub      *talk_together_app.Hub
}

func NewHandler(services *service.Service, hub *talk_together_app.Hub) *Handler {
	return &Handler{services: services, hub: hub}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		room := api.Group("/room")
		{
			room.POST("/", h.createRoom)
			room.GET("/", h.getAllRooms)
			room.GET("/search", h.searchRoomByName)
			room.GET("/:id", h.getRoomById)
			room.PUT("/:id", h.updateRoom)
			room.DELETE("/:id", h.deleteRoom)
			room.GET("/:id/ws", h.handleWebSocket)
		}

		messages := api.Group("/messages")
		{
			messages.GET("/room/:room_id", h.getRoomMessages)
			messages.POST("/", h.sendMessage)
			messages.DELETE("/:id", h.deleteMessage)
			messages.PATCH("/:id", h.updateMessage)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8000")

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now(),
		})
	})

	return router
}
