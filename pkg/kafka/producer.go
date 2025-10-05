package kafka

import (
	"github.com/IBM/sarama"
	"github.com/firstproject/talk-together-app/model"
	"github.com/goccy/go-json"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer, topic: topic}, nil
}

func (p *Producer) SendMessage(message interface{}, m model.Message) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(jsonMessage),
	}

	_, _, err = p.producer.SendMessage(msg)
	return err
}
