package hub

import (
	"github.com/firstproject/talk-together-app/model"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

// Представляет команту с клиентами и обеспечивает потокобезопасность
type RoomEntry struct {
	Room    *model.Room
	Clients map[int]*model.Client
	mu      sync.RWMutex
}

// Управляет WebSocket соединениями, комнатами и рассылкой сообщений
type Hub struct {
	Rooms      map[int]*RoomEntry
	mu         sync.RWMutex
	Upgrader   websocket.Upgrader
	Register   chan *model.Client
	Unregister chan *model.Client
	Broadcast  chan *model.Message
}

// Создает и инициализирует новый экземпляр Hub
func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[int]*RoomEntry),
		Register:   make(chan *model.Client),
		Unregister: make(chan *model.Client),
		Broadcast:  make(chan *model.Message),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// Запускает основной цикл обработки событий Hub
// Обрабатывает: регистрацию, отмену регистрации, broadcast сообщений
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.unregisterClient(client)
		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

// Создает новую комнату в Hub
// Если комната с таким Id существует , ничего не делает
func (h *Hub) CreateRoom(room *model.Room) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.Rooms[room.Id]; !exists {
		h.Rooms[room.Id] = &RoomEntry{
			Room:    room,
			Clients: make(map[int]*model.Client),
		}
	}
}

// Возвращает комнату по Id
// Возвращает nil если комната не найдена
func (h *Hub) GetRoom(roomId int) *model.Room {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if roomEntry, exists := h.Rooms[roomId]; exists {
		return roomEntry.Room
	}
	return nil
}

// Удаляет комнату и отключает всех ее клиентов
// Используется при удалении комнаты из системы
func (h *Hub) RemoveRoom(roomId int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if roomEntry, exists := h.Rooms[roomId]; exists {
		roomEntry.mu.Lock()
		for _, client := range roomEntry.Clients {
			close(client.Send)
			client.Conn.Close()
		}
		roomEntry.mu.Unlock()

		delete(h.Rooms, roomId)

	}

}

// Добавляет клиента в комнату
// Если комнаты не существует - создает ее
// Внутренний метод, вызывается из Run()
func (h *Hub) registerClient(client *model.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.Rooms[client.Room]; !exists {
		h.Rooms[client.Room] = &RoomEntry{
			Clients: make(map[int]*model.Client),
		}
	}
	h.Rooms[client.Room].mu.Lock()
	h.Rooms[client.Room].Clients[client.Id] = client
	h.Rooms[client.Room].mu.Unlock()
}

// Удаляет клиента из комнаты и закрывает соединение
// Если комната пустеет - автоматически удаляется
// Внутренний метод, вызывается из Run()
func (h *Hub) unregisterClient(client *model.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, exists := h.Rooms[client.Room]; exists {
		room.mu.Lock()
		delete(room.Clients, client.Id)
		close(client.Send)
		room.mu.Unlock()

		if len(room.Clients) == 0 {
			delete(h.Rooms, client.Room)
		}
	}
}

func (h *Hub) broadcastMessage(message *model.Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, exists := h.Rooms[message.Room]; exists {
		room.mu.RLock()
		defer room.mu.RUnlock()

		for _, client := range room.Clients {
			select {
			case client.Send <- []byte(message.Content):
			default:
				close(client.Send)
				delete(room.Clients, client.Id)
			}
		}
	}
}

func (h *Hub) HasRoom(roomId int) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, exists := h.Rooms[roomId]
	return exists
}

func (h *Hub) GetRoomClientsCount(roomId int) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if roomEntry, exists := h.Rooms[roomId]; exists {
		roomEntry.mu.Lock()
		defer roomEntry.mu.Unlock()
		return len(roomEntry.Clients)
	}
	return 0
}
