package kafka

import (
	"github.com/IBM/sarama"
	talk_together_app "github.com/firstproject/talk-together-app/hub"
	"github.com/firstproject/talk-together-app/model"
	"github.com/goccy/go-json"
	"log"
)

type Consumer struct {
	consumer sarama.Consumer
	topic    string
	hub      *talk_together_app.Hub
}

func NewKafkaConsumer(brokers []string, topic string, hub *talk_together_app.Hub) (*Consumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer, topic, hub}, nil
}

func (c *Consumer) Start() {
	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal(err)
	}

	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var message model.Message
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				c.hub.Broadcast <- &message
			}
		}
	}
}
