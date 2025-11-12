package services

// import (
// 	"context"
// 	"fmt"

// 	"github.com/TranVinhHien/ecom_payment-service/kafka"
// )

// // ============================================================================
// // HƯỚNG DẪN SỬ DỤNG KAFKA TRONG SERVICE
// // ============================================================================

// // 1. CÁCH GỬI MESSAGE ĐƠN GIẢN
// // ---------------------------
// // Sử dụng Producer để gửi message trực tiếp
// func (s *service) ExampleSendSimpleMessage(ctx context.Context) error {
// 	producer := s.kafkaClient.Producer()

// 	// Dữ liệu cần gửi
// 	messageData := map[string]interface{}{
// 		"user_id": "123456",
// 		"action":  "payment_initiated",
// 		"amount":  100000,
// 	}

// 	// Gửi message tới topic
// 	// - topic: tên topic muốn gửi
// 	// - key: key để phân vùng (partition), thường dùng user_id hoặc order_id
// 	// - value: dữ liệu message (sẽ được tự động convert sang JSON)
// 	err := producer.SendMessage(
// 		kafka.TopicPaymentCreated,  // Topic
// 		"user-123456",               // Key (dùng để phân vùng)
// 		messageData,                 // Data
// 	)

// 	if err != nil {
// 		return fmt.Errorf("lỗi gửi message: %w", err)
// 	}

// 	return nil
// }

// // 2. CÁCH GỬI MESSAGE VỚI HEADERS
// // -------------------------------
// // Headers thường dùng để truyền metadata như trace_id, user_id, request_id
// func (s *service) ExampleSendMessageWithHeaders(ctx context.Context) error {
// 	producer := s.kafkaClient.Producer()

// 	messageData := map[string]interface{}{
// 		"transaction_id": "TXN123",
// 		"status":         "completed",
// 	}

// 	// Tạo headers
// 	headers := map[string]string{
// 		"user_id":       "user-123",
// 		"trace_id":      "trace-abc-xyz",
// 		"service":       "payment-service",
// 		"content_type":  "application/json",
// 	}

// 	err := producer.SendMessageWithHeaders(
// 		kafka.TopicTransactionCompleted,
// 		"TXN123",
// 		messageData,
// 		headers,
// 	)

// 	return err
// }

// // 3. CÁCH GỬI NHIỀU MESSAGE (BATCH)
// // ---------------------------------
// // Dùng khi cần gửi nhiều message cùng lúc, hiệu quả hơn gửi từng message
// func (s *service) ExampleSendBatchMessages(ctx context.Context) error {
// 	producer := s.kafkaClient.Producer()

// 	// Tạo danh sách messages
// 	messages := []kafka.ProducerMessage{
// 		{
// 			Key: "order-001",
// 			Value: map[string]interface{}{
// 				"order_id": "order-001",
// 				"status":   "pending",
// 			},
// 		},
// 		{
// 			Key: "order-002",
// 			Value: map[string]interface{}{
// 				"order_id": "order-002",
// 				"status":   "pending",
// 			},
// 			Headers: map[string]string{
// 				"priority": "high",
// 			},
// 		},
// 	}

// 	// Gửi batch
// 	err := producer.SendMessageBatch(kafka.TopicOrderPaymentReceived, messages)
// 	return err
// }

// // 4. CÁCH GỬI MESSAGE BẤT ĐỒNG BỘ (ASYNC)
// // ----------------------------------------
// // Dùng khi không cần đợi confirmation, chỉ fire and forget
// func (s *service) ExampleSendAsyncMessage(ctx context.Context) error {
// 	producer := s.kafkaClient.Producer()

// 	// Gửi async - không block thread hiện tại
// 	err := producer.SendMessageAsync(
// 		kafka.TopicNotificationEmail,
// 		"user-123",
// 		map[string]interface{}{
// 			"to":      "user@example.com",
// 			"subject": "Payment confirmation",
// 			"body":    "Your payment was successful",
// 		},
// 	)

// 	// Function này return ngay lập tức
// 	return err
// }

// // 5. CÁCH SỬ DỤNG EVENT PUBLISHER (RECOMMENDED)
// // ---------------------------------------------
// // Event Publisher cung cấp các methods có sẵn cho từng loại event
// func (s *service) ExampleUseEventPublisher(ctx context.Context) error {
// 	// A. Publish Payment Created Event
// 	err := s.eventPublisher.PublishPaymentCreated(
// 		ctx,
// 		"payment-123",      // Payment ID
// 		"user-456",         // User ID
// 		100000.0,           // Amount
// 		map[string]interface{}{ // Additional data
// 			"payment_method": "MOMO",
// 			"order_id":       "order-789",
// 		},
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	// B. Publish Payment Completed Event
// 	err = s.eventPublisher.PublishPaymentCompleted(
// 		ctx,
// 		"payment-123",
// 		"transaction-xyz",
// 		map[string]interface{}{
// 			"gateway_response": "SUCCESS",
// 		},
// 	)

// 	// C. Publish Payment Failed Event
// 	err = s.eventPublisher.PublishPaymentFailed(
// 		ctx,
// 		"payment-123",
// 		"Insufficient balance",
// 		nil,
// 	)

// 	return err
// }

// // 6. CÁCH GỬI MESSAGE TRONG GOROUTINE (NON-BLOCKING)
// // --------------------------------------------------
// // Dùng khi muốn gửi message mà không làm chậm logic chính
// func (s *service) ExampleSendInGoroutine(ctx context.Context, orderID string, amount float64) {
// 	// Logic chính của service...
// 	fmt.Println("Đang xử lý đơn hàng...")

// 	// Gửi event trong goroutine để không block
// 	go func() {
// 		eventData := map[string]interface{}{
// 			"order_id": orderID,
// 			"amount":   amount,
// 			"status":   "processing",
// 		}

// 		if err := s.eventPublisher.PublishPaymentCreated(context.Background(), orderID, "user-123", amount, eventData); err != nil {
// 			// Log lỗi nhưng không ảnh hưởng tới flow chính
// 			fmt.Println("⚠️ Lỗi gửi Kafka event:", err)
// 		}
// 	}()

// 	// Tiếp tục xử lý logic...
// 	fmt.Println("Đơn hàng đã được xử lý")
// }

// // 7. BEST PRACTICES KHI GỬI MESSAGE
// // ----------------------------------
// func (s *service) ExampleBestPractices(ctx context.Context) error {
// 	// ✅ GOOD: Sử dụng context để timeout
// 	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	// ✅ GOOD: Sử dụng key có ý nghĩa để phân vùng tốt hơn
// 	// Nên dùng: user_id, order_id, transaction_id
// 	orderID := "order-123"

// 	// ✅ GOOD: Gửi trong goroutine để không block
// 	go func() {
// 		err := s.eventPublisher.PublishPaymentCreated(
// 			ctxWithTimeout,
// 			"payment-123",
// 			"user-456",
// 			100000,
// 			map[string]interface{}{
// 				"order_id": orderID,
// 			},
// 		)

// 		// ✅ GOOD: Luôn log error
// 		if err != nil {
// 			fmt.Printf("❌ Lỗi gửi Kafka event: %v\n", err)
// 			// Có thể lưu vào retry queue hoặc dead letter queue
// 		}
// 	}()

// 	return nil
// }

// // 8. CÁC TOPIC CÓ SẴN
// // -------------------
// // Tất cả topics được định nghĩa trong kafka/topics.go:
// //
// // Payment topics:
// // - kafka.TopicPaymentCreated
// // - kafka.TopicPaymentCompleted
// // - kafka.TopicPaymentFailed
// // - kafka.TopicPaymentRefunded
// //
// // Transaction topics:
// // - kafka.TopicTransactionCreated
// // - kafka.TopicTransactionCompleted
// // - kafka.TopicTransactionTimeout
// // - kafka.TopicTransactionCancelled
// //
// // Order topics:
// // - kafka.TopicOrderPaymentReceived
// // - kafka.TopicOrderPaymentFailed
// //
// // Notification topics:
// // - kafka.TopicNotificationEmail
// // - kafka.TopicNotificationSMS
// // - kafka.TopicNotificationPush

// // ============================================================================
// // VÍ DỤ THỰC TÉ: Tích hợp vào Payment Flow
// // ============================================================================

// func (s *service) ExampleRealPaymentFlow(ctx context.Context, userID, orderID string, amount float64) error {
// 	// 1. Tạo payment record trong DB
// 	paymentID := "payment-" + orderID

// 	// 2. Gửi event Payment Created (async để không block)
// 	go func() {
// 		err := s.eventPublisher.PublishPaymentCreated(
// 			context.Background(),
// 			paymentID,
// 			userID,
// 			amount,
// 			map[string]interface{}{
// 				"order_id":       orderID,
// 				"payment_method": "MOMO",
// 				"currency":       "VND",
// 			},
// 		)
// 		if err != nil {
// 			fmt.Println("❌ Lỗi gửi PaymentCreated event:", err)
// 		}
// 	}()

// 	// 3. Xử lý thanh toán...
// 	success := true // giả sử thành công

// 	if success {
// 		// 4. Gửi event Payment Completed
// 		go func() {
// 			err := s.eventPublisher.PublishPaymentCompleted(
// 				context.Background(),
// 				paymentID,
// 				"txn-123",
// 				map[string]interface{}{
// 					"order_id":    orderID,
// 					"gateway":     "MOMO",
// 					"completed_at": time.Now(),
// 				},
// 			)
// 			if err != nil {
// 				fmt.Println("❌ Lỗi gửi PaymentCompleted event:", err)
// 			}
// 		}()
// 	} else {
// 		// 5. Gửi event Payment Failed
// 		go func() {
// 			err := s.eventPublisher.PublishPaymentFailed(
// 				context.Background(),
// 				paymentID,
// 				"Gateway timeout",
// 				nil,
// 			)
// 			if err != nil {
// 				fmt.Println("❌ Lỗi gửi PaymentFailed event:", err)
// 			}
// 		}()
// 	}

// 	return nil
// }
