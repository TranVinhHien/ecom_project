package redis_db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/TranVinhHien/ecom_order_service/services"
	modelServices "github.com/TranVinhHien/ecom_order_service/services/entity"

	"github.com/go-co-op/gocron/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
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
func (s *RedisDB) addScoreMember(ctx context.Context, zsetKey, token string, expiry float64) error {
	_, err := s.client.ZAdd(ctx, zsetKey, redis.Z{
		Score: expiry,
		//Score:  float64(time.Now().Add(time.Second * 15).Unix()),
		Member: token,
	}).Result()
	if err != nil {
		return fmt.Errorf("error when add new item %s: to key storeMember:%s ,err: %v", token, zsetKey, err)
	}
	return nil
}

// // Hàm thêm token vào ZSET với thời gian hết hạn
// func (s *RedisDB) AddOrderOnline(ctx context.Context, user_id string, payload modelServices.CombinedDataPayLoadMoMo, duration time.Duration) error {
// 	jsonValue, err := json.Marshal(payload)
// 	if err != nil {
// 		return fmt.Errorf("lỗi khi chuyển CombinedDataPayLoadMoMo sang JSON: %v", err)
// 	}
// 	value := string(jsonValue)
// 	errRedis := s.client.Set(ctx, OrderOnline+user_id+"_"+payload.OrderTX.OrderID, value, duration)
// 	if errRedis != nil {
// 		return fmt.Errorf("error AddOrderOnline when add new item %s: to key :%s ,err: %v", user_id, value, errRedis)
// 	}
// 	return nil
// }

// func (s *RedisDB) GetOrderOnline(ctx context.Context, user_id string) (payload *modelServices.CombinedDataPayLoadMoMo, err error) {
// 	// value, errRedis := s.client.Get(ctx, OrderOnline+user_id).Result()
// 	// if errRedis != nil {
// 	// 	return nil, fmt.Errorf("error GetOrderOnline  key :%s ,err: %v", user_id, errRedis)
// 	// }
// 	// if value == "" {
// 	// 	return nil, nil
// 	// }
// 	// err = json.Unmarshal([]byte(value), &payload)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("lỗi chuyển đổi JSON: %v", err)
// 	// }
// 	// return
// 	pattern := fmt.Sprintf("%s%s_*", OrderOnline, user_id)

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

// scan item expired and remove it
func (s *RedisDB) removeExpired(ctx context.Context, zsetKey string) error {
	now := float64(time.Now().Unix())
	// Lấy danh sách token hết hạn
	expiredTokens, err := s.client.ZRangeByScore(ctx, zsetKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%f", now),
	}).Result()
	if err != nil {
		return fmt.Errorf("error when check token: %v", err)
	}
	if len(expiredTokens) == 0 {
		return nil
	}
	// Xóa các token hết hạn
	_, err = s.client.ZRem(ctx, zsetKey, expiredTokens).Result()
	if err != nil {

		return fmt.Errorf("error when remove token: %v", err)
	}
	log.Info().Msg(fmt.Sprintf("remove %v token", len(expiredTokens)))
	return nil
}

// check token Valid
func (s *RedisDB) isExists(ctx context.Context, zsetKey, token string) bool {
	// Dùng ZSCORE để kiểm tra token
	now := float64(time.Now().Unix()) // Lấy thời gian hiện tại

	score, err := s.client.ZScore(ctx, zsetKey, token).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		return false
	}
	if now >= score {
		return false
	}
	// Token tồn tại
	return true
}
func (s *RedisDB) AddTokenToBlackList(ctx context.Context, token string, exprid float64) error {
	return s.addScoreMember(ctx, BLACK_LIST, token, exprid)
}
func (s *RedisDB) CheckExistsFromBlackList(ctx context.Context, token string, exprid float64) bool {
	return s.isExists(ctx, BLACK_LIST, token)
}

// auto remove token expired
func (s *RedisDB) RemoveTokenExp(zsetKey string) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		fmt.Print(err)
	}
	j, err := scheduler.NewJob(
		gocron.DurationJob(
			24*time.Hour, //
		),
		gocron.NewTask(
			func() {
				s.removeExpired(context.Background(), zsetKey)
			},
		),
	)
	if err != nil {
		// handle error
	}
	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	log.Info().Msg("Started job")
	scheduler.Start()
}

func (s *RedisDB) AddCategories(ctx context.Context, cates []modelServices.Categorys) error {

	// tao 2 dataset truoc
	dataMap := make(map[string]string)
	childMap := make(map[string][]string)
	for _, cat := range cates {
		catJson, err := json.Marshal(cat)
		if err != nil {
			return err
		}
		dataMap[cat.CategoryID] = string(catJson)
		parentID := cat.CategoryID
		if cat.Parent.Valid {
			parentID = cat.Parent.Data
		}
		if parentID != cat.CategoryID {
			childMap[parentID] = append(childMap[parentID], cat.CategoryID)
		}
	}
	// Lưu dataMap lên Redis (HMSET)
	if err := s.client.HSet(ctx, CategoryDataKey, dataMap).Err(); err != nil {
		return fmt.Errorf("failed to save category data: %w", err)
	}
	// Lưu childrenMap lên Redis
	for parentID, children := range childMap {
		childrenJson, err := json.Marshal(children)
		if err != nil {
			return fmt.Errorf("failed to marshal children for parent %s: %w", parentID, err)
		}
		if err := s.client.HSet(ctx, CategoryChildrenKey, parentID, childrenJson).Err(); err != nil {
			return fmt.Errorf("failed to save children: %w", err)
		}
	}
	return nil
}
func (s *RedisDB) RemoveCategories(ctx context.Context) error {
	err := s.client.Del(ctx, CategoryChildrenKey).Err()
	if err != nil {
		return err
	}
	return s.client.Del(ctx, CategoryDataKey).Err()

}
func (s *RedisDB) GetCategoryTree(ctx context.Context, rootID string) ([]modelServices.Categorys, error) {
	// Lấy toàn bộ dữ liệu danh mục từ Redis
	dataMap, err := s.client.HGetAll(ctx, CategoryDataKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get category data: %w", err)
	}

	// Lấy toàn bộ dữ liệu con từ Redis
	childMap, err := s.client.HGetAll(ctx, CategoryChildrenKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get category children: %w", err)
	}

	// Parse dữ liệu thành map các danh mục
	categories := make(map[string]*modelServices.Categorys)
	for _, jsonStr := range dataMap {
		var cat modelServices.Categorys
		if err := json.Unmarshal([]byte(jsonStr), &cat); err != nil {
			return nil, fmt.Errorf("failed to unmarshal category: %w", err)
		}
		categories[cat.CategoryID] = &cat
	}

	// Hàm đệ quy để xây dựng cây con với độ sâu tối đa 3 lớp
	var buildTree func(parentID string, currentDepth int) ([]modelServices.Categorys, error)
	buildTree = func(parentID string, currentDepth int) ([]modelServices.Categorys, error) {
		childrenJson, exists := childMap[parentID]
		if !exists || currentDepth >= 3 { // Giới hạn độ sâu 3 lớp
			return nil, nil
		}

		var childrenIDs []string
		if err := json.Unmarshal([]byte(childrenJson), &childrenIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal children for parent %s: %w", parentID, err)
		}

		var childCategories []modelServices.Categorys
		for _, childID := range childrenIDs {
			if child, exists := categories[childID]; exists {
				// Tạo bản sao của category để tránh tham chiếu trực tiếp
				childCopy := *child

				// Đệ quy xây dựng cây con với độ sâu tăng lên
				grandchildren, err := buildTree(childID, currentDepth+1)
				if err != nil {
					return nil, err
				}

				if len(grandchildren) > 0 {
					childCopy.Childs = modelServices.Narg[[]modelServices.Categorys]{
						Data:  grandchildren,
						Valid: true,
					}
				}

				childCategories = append(childCategories, childCopy)
			}
		}

		return childCategories, nil
	}

	// Nếu có rootID, chỉ lấy danh mục đó và các con của nó (tối đa 3 lớp)
	if rootID != "" {
		rootCat, exists := categories[rootID]
		if !exists {
			return nil, fmt.Errorf("category with ID %s not found", rootID)
		}

		// Tạo bản sao của root category
		rootCopy := *rootCat

		// Xây dựng cây con với độ sâu bắt đầu từ 1 (vì root là lớp 0)
		children, err := buildTree(rootID, 1)
		if err != nil {
			return nil, err
		}

		if len(children) > 0 {
			rootCopy.Childs = modelServices.Narg[[]modelServices.Categorys]{
				Data:  children,
				Valid: true,
			}
		}

		return []modelServices.Categorys{rootCopy}, nil
	}

	// Nếu không có rootID, lấy toàn bộ danh mục gốc và xây dựng cây con
	var result []modelServices.Categorys
	for _, cat := range categories {
		if !cat.Parent.Valid || categories[cat.Parent.Data] == nil {
			// Tạo bản sao của category gốc
			rootCopy := *cat

			// Xây dựng cây con với độ sâu bắt đầu từ 1
			children, err := buildTree(cat.CategoryID, 1)
			if err != nil {
				return nil, err
			}

			if len(children) > 0 {
				rootCopy.Childs = modelServices.Narg[[]modelServices.Categorys]{
					Data:  children,
					Valid: true,
				}
			}

			result = append(result, rootCopy)
		}
	}

	return result, nil
}

func (s *RedisDB) DeleteOrderOnline(ctx context.Context, orderID string) error {
	// Giả sử hằng OrderOnline có giá trị "OrderOnline:"
	// Pattern sẽ là "OrderOnline:*_<orderID>"
	pattern := fmt.Sprintf("%s*_%s", OrderOnline, orderID)

	// Lấy các key phù hợp. Nếu dữ liệu của bạn không quá lớn, có thể dùng KEYS
	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("lỗi quét key với pattern %s: %w", pattern, err)
	}

	if len(keys) == 0 {
		// Không tìm thấy key nào, có thể trả về lỗi hoặc đơn giản là không làm gì
		return fmt.Errorf("không tìm thấy key nào cho order id %s", orderID)
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

// Bắt đầu lắng nghe sự kiện hết hạn key
func (h *RedisDB) StartExpirationListenerOrderOnline(cb func(ctx context.Context, orderID string)) {
	ctx := context.Background()
	fmt.Println("Starting expiration listener for OrderOnline keys...")
	pubsub := h.client.Subscribe(context.Background(), "__keyevent@0__:expired")
	// Chạy listener trong goroutine
	go func() {
		defer pubsub.Close()
		channel := pubsub.Channel()
		for msg := range channel {

			expiredKey := msg.Payload
			fmt.Println("Key expired:", expiredKey)
			if strings.HasPrefix(expiredKey, OrderOnline) {
				fmt.Println("OrderOnline:", OrderOnline)
				strings := strings.Split(expiredKey, "_")
				if len(strings) > 0 {
					fmt.Println("Starting expiration listener for OrderOnline keys...", strings[1])
					cb(ctx, strings[1])

				}
			}
		}
	}()
}
