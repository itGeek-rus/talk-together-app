package hub

import (
	"github.com/firstproject/talk-together-app/model"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockConn struct {
	mock.Mock
	writeChan chan []byte
}

func (m *MockConn) ReadMessage() (int, []byte, error) {
	args := m.Called()
	return args.Int(0), args.Get(1).([]byte), args.Error(1)
}

func (m *MockConn) WriteMessage(messageType int, data []byte) error {
	m.Called(messageType, data)
	if m.writeChan != nil {
		m.writeChan <- data
	}
	return nil
}

func (m *MockConn) Close() error {
	m.Called()
	return nil
}

func TestNewHub(t *testing.T) {
	hub := NewHub()
	assert.NotNil(t, hub)
	assert.NotNil(t, hub.Rooms)
	assert.NotNil(t, hub.Register)
	assert.NotNil(t, hub.Unregister)
	assert.NotNil(t, hub.Broadcast)
}

func TestHub_RegisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer close(hub.Register)

	client := &model.Client{
		Id:   1,
		Conn: &websocket.Conn{},
		Room: 1,
		User: 1,
		Send: make(chan []byte, 10),
	}

	hub.Register <- client

	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	defer hub.mu.RUnlock()

	assert.Contains(t, hub.Rooms, 1)
	assert.Contains(t, hub.Rooms[1].Clients, 1)
}

func TestHub_BroadcastMessage(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &model.Client{
		Id:   1,
		Conn: &websocket.Conn{},
		Room: 1,
		User: 1,
		Send: make(chan []byte, 10),
	}

	hub.Register <- client
	time.Sleep(10 * time.Millisecond)

	message := &model.Message{
		Id:      1,
		Room:    1,
		User:    1,
		Content: "test message",
	}

	hub.Broadcast <- message
	time.Sleep(10 * time.Millisecond)

	select {
	case recivedMsg := <-client.Send:
		assert.Equal(t, "test message", string(recivedMsg))
	case <-time.After(100 * time.Millisecond):
		t.Error("Message was not send to client")
	}

}

func TestHub_UnregisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &model.Client{
		Id:   1,
		Conn: &websocket.Conn{},
		Room: 1,
		User: 1,
		Send: make(chan []byte, 10),
	}

	hub.Register <- client
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	assert.Contains(t, hub.Rooms, 1)
	assert.Contains(t, hub.Rooms[1].Clients, 1)
	hub.mu.RUnlock()

	hub.Unregister <- client
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	defer hub.mu.RUnlock()

	if room, exists := hub.Rooms[1]; exists {
		room.mu.RLock()
		defer room.mu.RUnlock()
		assert.NotContains(t, room.Clients, 1)
	} else {
		t.Log("Room was automatically remove after last client unregistered")
	}

}

func TestHub_CreateRoom(t *testing.T) {
	hub := NewHub()

	room := &model.Room{
		Id:          1,
		Name:        "Test Room",
		Description: "Test Description",
		CreatedBy:   1,
	}

	hub.CreateRoom(room)

	hub.mu.RLock()
	defer hub.mu.RUnlock()

	assert.Contains(t, hub.Rooms, 1)
	assert.Equal(t, room, hub.Rooms[1].Room)
}
