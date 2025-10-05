package handler

import (
	talk_together_app "github.com/firstproject/talk-together-app/hub"
	"github.com/firstproject/talk-together-app/model"
	"github.com/firstproject/talk-together-app/pkg/middleware/monitoring"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var hub *talk_together_app.Hub

func init() {
	hub = talk_together_app.NewHub()
}

// @Summary WebSocket для чат-комнаты
// @Description Real-time WebSocket соединение для обмена сообщениями в комнате
// @Tags chat
// @Param roomId path int true "ID комнаты" example(1)
// @Param userId path int true "ID пользователя" example(1)
// @Success 101 "Switching Protocols"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /ws/rooms/{roomId}/users/{userId} [get]
func (h *Handler) handleWebSocket(c *gin.Context) {
	roomId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.services.Client.AddClientToRoom(roomId, userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	conn, err := hub.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	monitoring.IncrementWebSocketConnections()
	defer monitoring.DecrementWebSocketConnections()

	client := &model.Client{
		Id:   userId,
		Conn: conn,
		Room: roomId,
		User: userId,
		Send: make(chan []byte, 256),
	}

	hub.Register <- client

	go h.readPump(client)
	go h.writePump(client)

}

func (h *Handler) readPump(client *model.Client) {
	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		msgId, err := h.services.Message.CreateMessage(client.Room, client.Id, string(message))
		if err != nil {
			continue
		}

		hub.Broadcast <- &model.Message{
			Id:      msgId,
			Room:    client.Room,
			User:    client.User,
			Content: string(message),
		}
	}

}

func (h *Handler) writePump(client *model.Client) {
	defer client.Conn.Close()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}

}
