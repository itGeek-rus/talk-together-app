package handler

import (
	"github.com/firstproject/talk-together-app/pkg/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Get room message
// @Security ApiKeyAuth
// @Tags messages
// @Description Get room messages
// @ID get-room-messages
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/messages/room_id/{room_id} [get]
func (h *Handler) getRoomMessages(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Param("room_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, " invalid room id")
		return
	}

	messages, err := h.services.GetRoomMessages(roomId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, messages)
}

// @Summary Send message
// @Security ApiKeyAuth
// @tags messages
// @Description Send message
// @ID send-message
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/messages [post]
func (h *Handler) sendMessage(c *gin.Context) {
	userId, err := middleware.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input struct {
		Room    int    `json:"room" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.CreateMessage(input.Room, userId, input.Content)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary DeleteMessage
// @Security ApiKeyAuth
// @Tags messages
// @Description Delete message by ID
// @ID delete-message
// @Accept json
// @Produce json
// @Param id path int true "Message ID"
// @Success 200 {object} StatusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/messages/{id} [delete]
func (h *Handler) deleteMessage(c *gin.Context) {
	userId, err := middleware.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	messageId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid message id")
		return
	}

	err = h.services.DeleteMessage(messageId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "message deleted",
	})

}

type updateMessageInput struct {
	Content string `json:"content" binding:"required"`
}

// @Summary Update message
// @Security ApiKeyAuth
// @Tags messages
// @Description Update message content
// @ID update-message
// @Accept json
// @Produce json
// @Param id path int true "Message ID"
// @Param input body updateMessageInput true "Update input"
// @Success 200 {object} StatusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/messages/{id} [patch]
func (h *Handler) updateMessage(c *gin.Context) {
	userId, err := middleware.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	messageId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid message id")
		return
	}

	var input updateMessageInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Content == "" {
		newErrorResponse(c, http.StatusBadRequest, "content cannot be empty")
		return
	}

	if len(input.Content) > 1000 {
		newErrorResponse(c, http.StatusBadRequest, "cannot too long")
		return
	}

	err = h.services.UpdateMessage(messageId, userId, input.Content)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "message updated",
	})

}
