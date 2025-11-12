package services

import (
	"context"
	"log"
	"time"

	db "github.com/TranVinhHien/ecom_payment-service/db/sqlc"
	// services "github.com/TranVinhHien/ecom_payment-service/services/entity"
)

type HandlePaymentFailedDataOrderService struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"` // Đây là order_id (cha)
	Reason        string `json:"reason"`
}

func (s *service) CheckTransactionTimeout(ctx context.Context) {
	log.Println("Running job: Checking for expired transactions...")

	// 1. Tính mốc thời gian hết hạn
	expiredBefore := time.Now().Add(-s.env.OrderDuration) // Giả sử thời gian hết hạn là 15 phút

	// 2. Lấy các giao dịch PENDING đã quá hạn
	expiredTxs, err := s.repository.GetExpiredPendingTransactions(ctx, expiredBefore)
	if err != nil {
		log.Printf("Error fetching expired transactions: %v", err)
		return
	}

	if len(expiredTxs) == 0 {
		log.Println("Job finished: No expired transactions found.")
		return
	}

	log.Printf("Found %d expired transactions. Processing...", len(expiredTxs))

	for _, tx := range expiredTxs {
		// 3. Cập nhật status thành FAILED
		err := s.repository.UpdateTransactionStatus(ctx, db.UpdateTransactionStatusParams{
			Status: "FAILED",
			ID:     tx.ID,
		}) //
		if err != nil {
			log.Printf("Error updating transaction %s to FAILED: %v", tx.ID, err)
			continue // Bỏ qua và xử lý cái tiếp theo
		}

		// 4. Chuẩn bị message cho RabbitMQ
		eventBody := map[string]interface{}{
			"transaction_id": tx.ID,
			"order_id":       tx.OrderID.String,
			"reason":         "Hết thời gian thanh toán",
		}

		// 5. Gửi sự kiện 'payment_failed'
		err = s.producer.PaymentFailed(
			ctx,
			"payment.failed", // Routing key
			eventBody,
		)
		if err != nil {
			log.Printf("Error publishing payment_failed event for order %s: %v", tx.OrderID.String, err)
		} else {
			log.Printf("Published payment_failed event for order %s", tx.OrderID.String)
		}

	}
	log.Println("Job finished processing expired transactions.")

}
