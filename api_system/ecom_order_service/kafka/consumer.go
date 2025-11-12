package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/TranVinhHien/ecom_order_service/services"
	entity "github.com/TranVinhHien/ecom_order_service/services/entity"
)

type HandlePaymentSucceedData struct {
	OrderID       string  `json:"order_id"`
	Amount        float64 `json:"amount"`
	TransactionID string  `json:"transaction_id"`
	PaymentMethod string  `json:"payment_method"`
	Message       string  `json:"message"`
}
type HandlePaymentFailedData struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"` // Đây là order_id (cha)
	Reason        string `json:"reason"`
}

// KafkaConsumerHandler là adapter, nó implement interface của Sarama
type KafkaConsumerHandler struct {
	service services.ServiceUseCase // "Service" chứa logic nghiệp vụ
	ready   chan bool
}

// NewKafkaConsumerHandler tạo một handler mới
func NewKafkaConsumerHandler(service services.ServiceUseCase) *KafkaConsumerHandler {
	return &KafkaConsumerHandler{
		service: service,
		ready:   make(chan bool),
	}
}

// Ready trả về channel để báo hiệu consumer đã sẵn sàng
func (h *KafkaConsumerHandler) Ready() <-chan bool {
	return h.ready
}
func (h *KafkaConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka consumer is setup and ready.")
	// Đóng channel 'ready' để báo cho main.go biết là đã sẵn sàng
	close(h.ready)
	return nil
}

// Cleanup được gọi khi session kết thúc
func (h *KafkaConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim là vòng lặp chính xử lý message
func (h *KafkaConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Lặp qua các message trong partition này
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Println("Message channel closed, exiting ConsumeClaim.")
				return nil
			}

			// Lấy context từ session để xử lý graceful shutdown
			ctx := session.Context()
			var err error

			// Phân loại message dựa trên Topic
			switch message.Topic {
			case TopicPaymentCompleted:
				var data HandlePaymentSucceedData
				err = json.Unmarshal(message.Value, &data)
				if err != nil {
					log.Printf("ERROR unmarshaling message (offset %d): %v. Will retry.", message.Offset, err)
					break
				}
				err = h.service.HandlePaymentSucceededEvent(ctx, entity.PaymentSucceededEvent{
					OrderID:       data.OrderID,
					TransactionID: data.TransactionID,
				})
				if err != nil {
					log.Printf("ERROR handling payment succeeded event (offset %d): %v. Will retry.", message.Offset, err)
				}

			case TopicPaymentFailed:
				var data HandlePaymentFailedData
				err = json.Unmarshal(message.Value, &data)
				if err != nil {
					log.Printf("ERROR unmarshaling message (offset %d): %v. Will retry.", message.Offset, err)
					break
				}
				err = h.service.HandlePaymentFailedEvent(ctx, entity.PaymentFailedEvent{
					OrderID:       data.OrderID,
					TransactionID: data.TransactionID,
					Reason:        data.Reason,
				})
				if err != nil {
					log.Printf("ERROR handling payment failed event (offset %d): %v. Will retry.", message.Offset, err)
				}
			default:
				log.Printf("WARN: Nhận được message từ topic lạ: %s", message.Topic)
			}

			// Xử lý kết quả
			if err != nil {
				// Nếu logic nghiệp vụ trả về lỗi, chúng ta KHÔNG commit.
				// Kafka sẽ tự động gửi lại message này.
				log.Printf("ERROR processing message (offset %d): %v. Will retry.", message.Offset, err)
			} else {
				// Xử lý thành công, commit message
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			// Thoát khi session bị hủy (ví dụ: service shutdown)
			return nil
		}
	}
}
