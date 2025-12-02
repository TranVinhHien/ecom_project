package redis_db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TranVinhHien/ecom_payment_service/services"
	modelServices "github.com/TranVinhHien/ecom_payment_service/services/entity"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	client *redis.Client
}
type RedisDBClient interface {
	// Save
}

func NewRedisDB(rdb *redis.Client) services.ServicesRedis {
	return &RedisDB{client: rdb}
}

// Hàm thêm token vào ZSET với thời gian hết hạn
func (s *RedisDB) AddTransactionOnline(ctx context.Context, user_id string, payload modelServices.CombinedDataPayLoadMoMo, duration time.Duration) error {
	jsonValue, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("lỗi khi chuyển CombinedDataPayLoadMoMo sang JSON: %v", err)
	}
	value := string(jsonValue)
	errRedis := s.client.Set(ctx, OrderOnline+user_id+"_"+payload.TransactionID, value, duration)
	if errRedis != nil {
		return fmt.Errorf("error AddTransactionOnline when add new item %s: to key :%s ,err: %v", user_id, value, errRedis)
	}
	return nil
}

func (s *RedisDB) GetTransactionOnline(ctx context.Context, user_id string) (payload *modelServices.CombinedDataPayLoadMoMo, err error) {
	// value, errRedis := s.client.Get(ctx, OrderOnline+user_id).Result()
	// if errRedis != nil {
	// 	return nil, fmt.Errorf("error GetTransactionOnline  key :%s ,err: %v", user_id, errRedis)
	// }
	// if value == "" {
	// 	return nil, nil
	// }
	// err = json.Unmarshal([]byte(value), &payload)
	// if err != nil {
	// 	return nil, fmt.Errorf("lỗi chuyển đổi JSON: %v", err)
	// }
	// return
	pattern := fmt.Sprintf("%s%s_*", OrderOnline, user_id)

	// Sử dụng SCAN để tìm key đầu tiên phù hợp
	var firstKey string
	var cursor uint64 = 0

	// Lấy 1 key mỗi lần quét
	keys, _, err := s.client.Scan(ctx, cursor, pattern, 1).Result()
	if err != nil {
		return nil, fmt.Errorf("lỗi quét key với pattern %s: %w", pattern, err)
	}
	if len(keys) == 0 {
		// Không tìm thấy key nào, trả về nil
		return nil, nil
	}
	firstKey = keys[0]

	// Lấy giá trị của key đầu tiên tìm thấy
	value, errRedis := s.client.Get(ctx, firstKey).Result()
	if errRedis != nil {
		return nil, fmt.Errorf("lỗi GetOrderOnline key: %s, err: %v", firstKey, errRedis)
	}
	if value == "" {
		return nil, nil
	}
	err = json.Unmarshal([]byte(value), &payload)
	if err != nil {
		return nil, fmt.Errorf("lỗi chuyển đổi JSON: %v", err)
	}
	return payload, nil
}

// func (s *RedisDB) GetTransactionOnlineWithIDTran(ctx context.Context, transactionID string) (payload *modelServices.CombinedDataPayLoadMoMo, err error) {

// 	pattern := fmt.Sprintf("*_%s", transactionID)

// 	// Sử dụng SCAN để tìm key đầu tiên phù hợp
// 	var firstKey string
// 	var cursor uint64 = 0

// 	// Lấy 1 key mỗi lần quét
// 	keys, _, err := s.client.Scan(ctx, cursor, pattern, 1).Result()
// 	if err != nil {
// 		return nil, fmt.Errorf("lỗi quét key với pattern %s: %w", pattern, err)
// 	}
// 	if len(keys) == 0 {
// 		// Không tìm thấy key nào, trả về nil
// 		return nil, nil
// 	}
// 	firstKey = keys[0]

// 	// Lấy giá trị của key đầu tiên tìm thấy
// 	value, errRedis := s.client.Get(ctx, firstKey).Result()
// 	if errRedis != nil {
// 		return nil, fmt.Errorf("lỗi GetOrderOnline key: %s, err: %v", firstKey, errRedis)
// 	}
// 	if value == "" {
// 		return nil, nil
// 	}
// 	err = json.Unmarshal([]byte(value), &payload)
// 	if err != nil {
// 		return nil, fmt.Errorf("lỗi chuyển đổi JSON: %v", err)
// 	}
// 	return payload, nil
// }

func (s *RedisDB) GetTransactionOnlineWithIDTran(ctx context.Context, transactionID string) (payload *modelServices.CombinedDataPayLoadMoMo, err error) {
	// Pattern này giả định key KẾT THÚC BẰNG _transactionID
	// Ví dụ: abc_12345
	// Xem Lỗi 2 nếu bạn không chắc chắn
	pattern := fmt.Sprintf("*_%s", transactionID)

	// In ra để debug (RẤT QUAN TRỌNG)
	// fmt.Println("Đang quét với pattern:", pattern)

	var firstKey string
	var cursor uint64 = 0
	var keys []string

	// BẮT BUỘC phải có vòng lặp
	for {
		// Yêu cầu quét, lấy về batch keys VÀ cursor tiếp theo
		keys, cursor, err = s.client.Scan(ctx, cursor, pattern, 1).Result() // COUNT = 1 để lấy 1 key mỗi lần
		if err != nil {
			return nil, fmt.Errorf("lỗi quét key với pattern %s: %w", pattern, err)
		}

		// Nếu batch này có key, lấy key đầu tiên và thoát lặp
		if len(keys) > 0 {
			firstKey = keys[0]
			break
		}

		// Nếu cursor = 0, nghĩa là đã quét HẾT toàn bộ database
		// mà không tìm thấy key nào
		if cursor == 0 {
			return nil, nil // Không tìm thấy
		}
	}

	// --- Phần code bên dưới giữ nguyên ---

	// Lấy giá trị của key đầu tiên tìm thấy
	value, errRedis := s.client.Get(ctx, firstKey).Result()
	if errRedis != nil {
		return nil, fmt.Errorf("lỗi GetOrderOnline key: %s, err: %v", firstKey, errRedis)
	}
	if value == "" {
		return nil, nil
	}
	err = json.Unmarshal([]byte(value), &payload)
	if err != nil {
		return nil, fmt.Errorf("lỗi chuyển đổi JSON: %v", err)
	}
	return payload, nil
}

func (s *RedisDB) DeleteTransactionOnline(ctx context.Context, transactionID string) error {
	// Giả sử hằng OrderOnline có giá trị "OrderOnline:"
	// Pattern sẽ là "OrderOnline:*_<transactionID>"
	pattern := fmt.Sprintf("%s*_%s", OrderOnline, transactionID)

	// Lấy các key phù hợp. Nếu dữ liệu của bạn không quá lớn, có thể dùng KEYS
	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("lỗi quét key với pattern %s: %w", pattern, err)
	}

	if len(keys) == 0 {
		// Không tìm thấy key nào, có thể trả về lỗi hoặc đơn giản là không làm gì
		return fmt.Errorf("không tìm thấy key nào cho transaction id %s", transactionID)
	}

	// Nếu chỉ muốn xóa key đầu tiên tìm thấy, có thể làm như sau:
	// err = s.client.Del(ctx, keys[0]).Err()
	// Nếu muốn xóa tất cả các key phù hợp, hãy truyền toàn bộ slice keys vào:
	err = s.client.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("lỗi xóa key(s): %w", err)
	}

	return nil
}
