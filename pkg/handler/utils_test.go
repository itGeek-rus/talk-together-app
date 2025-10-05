package handler

import (
	"github.com/firstproject/talk-together-app/model"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
)

func CreateTestContext(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(method, path, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	return c, w
}

type MockHub struct {
	RegisterChan   chan *model.Client
	UnregisterChan chan *model.Client
	BroadcastChan  chan *model.Client
}

func NewMockHub() *MockHub {
	return &MockHub{
		RegisterChan:   make(chan *model.Client, 10),
		UnregisterChan: make(chan *model.Client, 10),
		BroadcastChan:  make(chan *model.Client, 10),
	}
}
