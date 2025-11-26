package kafka

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
)

const (
	TopicPaymentCompleted = "payment.completed"
	TopicPaymentFailed    = "payment.failed"
)

// EventProducer là interface để các service của bạn sử dụng
type EventProducer interface {
	// Publish(ctx context.Context, topic string, key string, message []byte) error
	PaymentCompleted(ctx context.Context, key string, message map[string]interface{}) error
	PaymentFailed(ctx context.Context, key string, message map[string]interface{}) error
	Close() error
}

type kafkaProducer struct {
	producer sarama.SyncProducer
}

// NewProducer tạo một SyncProducer.
func NewProducer(brokers []string) (EventProducer, error) {
	config := GetSaramaConfig() // Lấy config chuẩn

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &kafkaProducer{producer: producer}, nil
}
func (p *kafkaProducer) PaymentCompleted(ctx context.Context, key string, message map[string]interface{}) error {
	// convert message to []byte
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return p.publish(ctx, TopicPaymentCompleted, key, messageBytes)
}
func (p *kafkaProducer) PaymentFailed(ctx context.Context, key string, message map[string]interface{}) error {
	// convert message to []byte
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return p.publish(ctx, TopicPaymentFailed, key, messageBytes)
}
func (p *kafkaProducer) publish(ctx context.Context, topic string, key string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(message),
	}

	// Gửi message và chờ xác nhận (vì là SyncProducer)
	_, _, err := p.producer.SendMessage(msg)
	return err
}

func (p *kafkaProducer) Close() error {
	return p.producer.Close()
}
