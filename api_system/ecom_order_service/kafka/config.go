package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

// GetSaramaConfig trả về một cấu hình sarama chuẩn cho producer/consumer
func GetSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()

	// Cấu hình chung
	config.Version = sarama.V2_8_0_0 // Đặt phiên bản Kafka broker của bạn
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second

	// Cấu hình Producer (quan trọng)
	config.Producer.Return.Successes = true // Bắt buộc để SyncProducer hoạt động
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll // An toàn nhất cho nghiệp vụ
	config.Producer.Retry.Max = 3

	// Cấu hình Consumer Group (quan trọng)
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.AutoCommit.Enable = true // Tự động commit
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	return config
}
