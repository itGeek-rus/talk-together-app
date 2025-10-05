package handler

import (
	"database/sql"
	"errors"
	"github.com/firstproject/talk-together-app/model"
	"github.com/firstproject/talk-together-app/pkg/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Create room
// @Security ApiKeyAuth
// @Tags room
// @Description Create room
// @ID create-room
// @Param input body model.Room true "Room info"
// @Success 200 {integer} 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Accept json
// @Produce json
// @Router /api/room [post]
func (h *Handler) createRoom(c *gin.Context) {
	userId, err := middleware.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input model.Room
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Room.CreateRoom(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

type getAllRoomsResponse struct {
	Data []model.Room `json:"data"`
}

// @Summary Get all rooms
// @Security ApiKeyAuth
// @Tags room
// @Description Get all rooms
// @ID get-all-rooms
// @Accept json
// @Produce json
// @Success 200 {object} getAllRoomsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/room [get]
func (h *Handler) getAllRooms(c *gin.Context) {
	userId, err := middleware.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	rooms, err := h.services.Room.GetAllRooms(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllRoomsResponse{
		Data: rooms,
	})
}

// @Summary Search rooms
// @Security ApiKeyAuth
// @Tags room
// @Description Search rooms
// @ID search-rooms
// @Accept json
// @Produce json
// @Success 200 {object} getAllRoomsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/room/search [get]
func (h *Handler) searchRoomByName(c *gin.Context) {
	_, err := middleware.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	searchQuery := c.Query("name")

	if searchQuery == "" {
		newErrorResponse(c, http.StatusBadRequest, "search query is empty")
		return
	}

	rooms, err := h.services.Room.SearchRoomByName(searchQuery)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllRoomsResponse{
		Data: rooms,
	})
}

type getRoomResponse struct {
	Data model.Room `json:"data"`
}

// @Summary Get room by id
// @Security ApiKeyAuth
// @Tags room
// @Description Get room by id
// @ID get-room-by-id
// @Accept json
// @Produce json
// @Success 200 {object} getRoomResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/room/:id [get]
func (h *Handler) getRoomById(c *gin.Context) {
	_, err := middleware.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid room id")
		return
	}

	room, err := h.services.Room.GetRoomById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newErrorResponse(c, http.StatusNotFound, "room not found")
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, getRoomResponse{
		Data: room,
	})
}

// @Summary Update room
// @Security ApiKeyAuth
// @Tags room
// @Description Update room
// @ID update-room
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/room/:id [put]
func (h *Handler) updateRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid room id")
		return
	}

	userId, err := middleware.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input model.UpdateRoomInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Room.UpdateRoom(userId, id, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newErrorResponse(c, http.StatusNotFound, "room not found")
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// @Summary Delete room
// @Security ApiKeyAuth
// @Tags room
// @Description Delete room
// @ID delete-room
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/room/:id [delete]
func (h *Handler) deleteRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	userId, err := middleware.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Room.DeleteRoom(id, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newErrorResponse(c, http.StatusNotFound, "room not found")
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
