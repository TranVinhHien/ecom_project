package main

// import gin
import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	config_assets "github.com/TranVinhHien/ecom_order_service/assets/config"
	"github.com/TranVinhHien/ecom_order_service/assets/token"
	"github.com/TranVinhHien/ecom_order_service/controllers"
	db "github.com/TranVinhHien/ecom_order_service/db/mysql"
	redis_db "github.com/TranVinhHien/ecom_order_service/db/redis"
	"github.com/TranVinhHien/ecom_order_service/kafka"
	"github.com/TranVinhHien/ecom_order_service/server"
	"github.com/TranVinhHien/ecom_order_service/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// create logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// read config file
	env, err := config_assets.LoadConfig(".")
	if err != nil {
		log.Err(err).Msg("Error read env:")
		return
	}

	// create connect to database
	conn, err := connectDBWithRetry(5, env.DBSource)
	if conn == nil {
		log.Err(err).Msg("Error when created connect to database")
		return
	}
	// close connection after gin stopped
	defer conn.Close()

	log.Info().Msg("Connect to database successfully")
	db := db.NewStore(conn)
	log.Info().Msg("Creating gin server...")
	// create jwt
	jwtMaker, err := token.NewJWTMaker(env.JWTSecret)
	if err != nil {
		log.Err(err).Msg("Error create JWTMaker")
		return
	}
	// create connect to redis
	rdb, err := connectDBRedisWithRetry(5, env.RedisAddress)
	if err != nil {
		log.Err(err).Msg("Error when created connect to redis")
		return
	}
	// create instance API server
	APIServer := server.NewAPIServices(env, time.Second*10)

	//setup redis Options
	redisdb := redis_db.NewRedisDB(rdb)
	// setup service
	services := services.NewService(db, jwtMaker, env, redisdb, APIServer)
	// setup controller
	controller := controllers.NewAPIController(services, jwtMaker)

	engine := gin.Default()
	engine.MaxMultipartMemory = 32 << 20 // 32 MB
	// engine.StaticFS("/.well-known", http.Dir("./assets/setup-mobile"))
	engine.GET("/.well-known/assetlinks.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.File("./assets/setup-mobile/assetlinks.json")
	})
	// config cors middleware
	config := cors.Config{
		AllowOrigins:     env.ClientIP,                                        // Chỉ cho phép localhost:3000
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Các method được phép
		AllowHeaders:     []string{"Content-Type", "Origin", "Authorization"}, // Các headers được phép
		ExposeHeaders:    []string{"Content-Length"},                          // Các headers trả về
		AllowCredentials: true,                                                // Cho phép cookies
	}

	v1 := engine.Group("/v1")
	v1.Use(cors.New(config))
	controller.SetUpRoute(v1)

	log.Info().Msg("Starting server on port " + env.HTTPServerAddress)

	// start jobs
	go redisdb.RemoveTokenExp(redis_db.BLACK_LIST)

	// Start Kafka consumer

	// --- 2. Khởi tạo Kafka Adapter (Handler) ---
	kafkaHandler := kafka.NewKafkaConsumerHandler(services)

	// --- 3. Khởi tạo Kafka Consumer Group ---
	brokers := []string{env.KafkaBrokers}
	topics := []string{kafka.TopicPaymentCompleted, kafka.TopicPaymentFailed}

	config_kafka := kafka.GetSaramaConfig() // Lấy config từ file producer của bạn
	consumerGroup, err := sarama.NewConsumerGroup(brokers, env.KafkaConsumerGroup, config_kafka)
	if err != nil {
		log.Err(err).Msg("Failed to create consumer group")
		return
	}
	defer consumerGroup.Close()
	go runConsumerGroup(consumerGroup, topics, kafkaHandler, services)
	engine.Run(env.HTTPServerAddress)

}
func connectDBRedisWithRetry(times int, redisAddress string) (*redis.Client, error) {
	var e error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2*time.Duration(times))
	defer cancel()
	for i := 1; i <= times; i++ {
		rdb := redis.NewClient(&redis.Options{
			Addr:     redisAddress,
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		_, err := rdb.Ping(ctx).Result()

		if err != nil {
			log.Err(err).Msg("Can't connect to redis")
		}
		// defer conn.Release()

		if err == nil {
			return rdb, nil
		}
		e = err
		time.Sleep(time.Second * 2)
	}
	return nil, e

}

//	func connectDBWithRetry(times int, dbConfig string) (*sql.DB, error) {
//		var e error
//		_, cancel := context.WithTimeout(context.Background(), time.Second*2*time.Duration(times))
//		defer cancel()
//		for i := 1; i <= times; i++ {
//			pool, err := sql.Open("mysql", dbConfig)
//			if err != nil {
//				log.Err(err).Msg("Can't create database pool")
//			}
//			err = pool.Ping()
//			if err != nil {
//				log.Err(err).Msg("Can't get connection to database pool")
//			}
//			// defer conn.Release()
//			pool.SetMaxOpenConns(10)                 // Số kết nối tối đa có thể mở
//			pool.SetMaxIdleConns(1)                  // Số kết nối có thể giữ mà không bị đóng
//			pool.SetConnMaxLifetime(5 * time.Minute) // Thời gian tối đa một kết nối có thể sống
//			if err == nil {
//				return pool, nil
//			}
//			e = err
//			time.Sleep(time.Second * 2)
//		}
//		return nil, e
//	}
func connectDBWithRetry(times int, dbConfig string) (*sql.DB, error) {
	var e error
	_, cancel := context.WithTimeout(context.Background(), time.Second*2*time.Duration(times))
	defer cancel()
	for i := 1; i <= times; i++ {
		// Thêm parseTime=true&loc=Asia%2FHo_Chi_Minh vào dbConfig
		if dbConfig[len(dbConfig)-1] == '/' {
			dbConfig += "?parseTime=true&loc=Asia%2FHo_Chi_Minh"
		} else {
			dbConfig += "&parseTime=true&loc=Asia%2FHo_Chi_Minh"
		}

		pool, err := sql.Open("mysql", dbConfig)
		if err != nil {
			log.Err(err).Msg("Can't create database pool")
		}
		err = pool.Ping()
		if err != nil {
			log.Err(err).Msg("Can't get connection to database pool")
		}

		pool.SetMaxOpenConns(10)
		pool.SetMaxIdleConns(1)
		pool.SetConnMaxLifetime(5 * time.Minute)

		if err == nil {
			// Set timezone cho session MySQL
			_, err = pool.Exec("SET time_zone = '+07:00'")
			if err != nil {
				log.Err(err).Msg("Can't set timezone")
			}
			return pool, nil
		}
		e = err
		time.Sleep(time.Second * 2)
	}
	return nil, e
}

func runConsumerGroup(consumerGroup sarama.ConsumerGroup, topics []string, kafkaHandler *kafka.KafkaConsumerHandler, services services.ServiceUseCase) {

	// --- 4. Bắt đầu chạy Consumer (liên tục) ---
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			// Consume sẽ block và chạy vòng lặp ConsumeClaim
			err := consumerGroup.Consume(ctx, topics, kafkaHandler)
			if err != nil {
				log.Print("Error from consumer: %w", err.Error())
				return
			}
			// Nếu context bị cancel (do shutdown), thoát vòng lặp
			if ctx.Err() != nil {
				return
			}
			// Reset 'ready' channel để chuẩn bị cho lần rebalance tiếp theo
			kafkaHandler = kafka.NewKafkaConsumerHandler(services)
		}
	}()

	// Chờ cho consumer sẵn sàng
	<-kafkaHandler.Ready()
	log.Print("Order Worker is up and running. Listening for messages...")

	// --- 5. Xử lý Graceful Shutdown ---
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	log.Print("Shutting down worker...")
	cancel()  // Gửi tín hiệu hủy cho context
	wg.Wait() // Chờ cho goroutine consumer kết thúc

	log.Print("Order Worker stopped.")

}
