package kafka

// import (
// 	"log"

// 	"github.com/IBM/sarama"
// 	"github.com/TranVinhHien/ecom_payment-service/services"
// )

// // KafkaConsumerHandler là adapter, nó implement interface của Sarama
// type KafkaConsumerHandler struct {
// 	service *services.ServiceUseCase // "Service" chứa logic nghiệp vụ
// 	ready   chan bool
// }

// // NewKafkaConsumerHandler tạo một handler mới
// func NewKafkaConsumerHandler(service *services.ServiceUseCase) *KafkaConsumerHandler {
// 	return &KafkaConsumerHandler{
// 		service: service,
// 		ready:   make(chan bool),
// 	}
// }

// // Ready trả về channel để báo hiệu consumer đã sẵn sàng
// func (h *KafkaConsumerHandler) Ready() <-chan bool {
// 	return h.ready
// }

// // Cleanup được gọi khi session kết thúc
// func (h *KafkaConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
// 	return nil
// }

// // ConsumeClaim là vòng lặp chính xử lý message
// func (h *KafkaConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
// 	// Lặp qua các message trong partition này
// 	for {
// 		select {
// 		case message, ok := <-claim.Messages():
// 			if !ok {
// 				log.Println("Message channel closed, exiting ConsumeClaim.")
// 				return nil
// 			}

// 			// Lấy context từ session để xử lý graceful shutdown
// 			ctx := session.Context()
// 			var err error

// 			// Phân loại message dựa trên Topic
// 			switch message.Topic {
// 			case TopicPaymentCompleted:
// 				err = h.service.HandlePaymentCompleted(ctx, string(message.Key), message.Value)
// 			case TopicPaymentFailed:
// 				err = h.service.HandlePaymentFailed(ctx, string(message.Key), message.Value)
// 			default:
// 				log.Printf("WARN: Nhận được message từ topic lạ: %s", message.Topic)
// 			}

// 			// Xử lý kết quả
// 			if err != nil {
// 				// Nếu logic nghiệp vụ trả về lỗi, chúng ta KHÔNG commit.
// 				// Kafka sẽ tự động gửi lại message này.
// 				log.Printf("ERROR processing message (offset %d): %v. Will retry.", message.Offset, err)
// 			} else {
// 				// Xử lý thành công, commit message
// 				session.MarkMessage(message, "")
// 			}

// 		case <-session.Context().Done():
// 			// Thoát khi session bị hủy (ví dụ: service shutdown)
// 			return nil
// 		}
// 	}
// }
