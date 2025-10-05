package redis

import (
	"fmt"
	"github.com/firstproject/talk-together-app/model"
	"github.com/firstproject/talk-together-app/pkg/middleware/monitoring"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

type Client struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(addr, password string, db int) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		logrus.Fatalf("Failed to connect to redis: %s", err.Error())
	}
	return &Client{client, context.Background()}
}

func (c *Client) CacheRoomMessages(roomId int, messages []model.Message) error {
	key := fmt.Sprintf("room:%d:messages", roomId)
	jsonData, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	return c.client.Set(c.ctx, key, jsonData, 10*time.Minute).Err()
}

func (c *Client) GetCachedRoomMessages(roomId int) ([]model.Message, error) {
	key := fmt.Sprintf("room:%d:messages", roomId)
	data, err := c.client.Get(c.ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var messages []model.Message
	err = json.Unmarshal(data, &messages)
	return messages, err
}

func (c *Client) Set(key string, value interface{}, expiration time.Duration) error {

	err := c.client.Set(c.ctx, key, value, expiration).Err()

	status := "success"
	if err != nil {
		status = "error"
	}

	monitoring.IncrementRedisOperations("set", status)

	return err
}

func (c *Client) Get(key string) (string, error) {
	result, err := c.client.Get(c.ctx, key).Result()

	status := "success"
	if err != nil {
		status = "error"
	}

	monitoring.IncrementRedisOperations("get", status)
	return result, err
}

func (c *Client) Del(key string) error {
	err := c.client.Del(c.ctx, key).Err()

	status := "success"
	if err != nil {
		status = "error"
	}

	monitoring.IncrementRedisOperations("del", status)
	return err
}

//Get, Del, HSet
