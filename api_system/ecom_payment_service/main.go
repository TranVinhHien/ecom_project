package main

// import gin
import (
	"context"
	"database/sql"
	"os"
	"time"

	config_assets "github.com/TranVinhHien/ecom_payment_service/assets/config"
	assets_jobs "github.com/TranVinhHien/ecom_payment_service/assets/jobs"
	"github.com/TranVinhHien/ecom_payment_service/assets/token"
	"github.com/TranVinhHien/ecom_payment_service/controllers"
	db "github.com/TranVinhHien/ecom_payment_service/db/mysql"
	redis_db "github.com/TranVinhHien/ecom_payment_service/db/redis"
	"github.com/TranVinhHien/ecom_payment_service/kafka"
	"github.com/TranVinhHien/ecom_payment_service/server"
	"github.com/TranVinhHien/ecom_payment_service/services"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	log.Info().Msg("Creating gin server...")
	// create jwt
	jwtMaker, err := token.NewJWTMaker(env.JWTSecret)
	if err != nil {
		log.Err(err).Msg("Error create JWTMaker")
		return
	}
	rdb, err := connectDBRedisWithRetry(5, env.RedisAddress)
	if err != nil {
		log.Err(err).Msg("Error when created connect to redis")
		return
	}
	//setup redis Options
	redisdb := redis_db.NewRedisDB(rdb)

	job, err := assets_jobs.NewJobScheduler()
	if err != nil {
		log.Err(err).Msg("Error when created job scheduler")
		return
	}

	// create connect to database
	conn, err := connectDBWithRetry(5, env.DBSource)
	if conn == nil {
		log.Err(err).Msg("Error when created connect to database")
		return
	}

	// create kafka producer
	producer, err := kafka.NewProducer([]string{env.KafkaBrokers})
	if err != nil {
		log.Err(err).Msg("Error when created kafka producer")
		return
	}
	log.Info().Msg("Connected to Kafka successfully")

	// close connection after gin stopped
	defer conn.Close()
	APIServer := server.NewAPIServices(env, time.Second*10)

	log.Info().Msg("Connect to database successfully")
	db := db.NewStore(conn)
	services := services.NewService(jwtMaker, env, redisdb, db, APIServer, producer)
	controller := controllers.NewAPIController(services, jwtMaker)

	engine := gin.Default()
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
	runJobs(services, job)
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

		pool.SetMaxOpenConns(60)
		pool.SetMaxIdleConns(60)
		pool.SetConnMaxLifetime(5 * time.Minute)
		pool.SetConnMaxIdleTime(2 * time.Minute)

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
func runJobs(s services.ServiceUseCase, jobScheduler *assets_jobs.JobScheduler) {
	ctx := context.Background()

	jobScheduler.NewJob(0, 1, 0, func() {
		// check transaction timeout
		s.CheckTransactionTimeout(ctx)
	})
	jobScheduler.Start()
}
