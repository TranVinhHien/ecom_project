package main

// import gin
import (
	"context"
	"database/sql"
	"os"
	"time"

	config_assets "github.com/TranVinhHien/ecom_analytics_service/assets/config"
	"github.com/TranVinhHien/ecom_analytics_service/assets/token"
	"github.com/TranVinhHien/ecom_analytics_service/controllers"
	db "github.com/TranVinhHien/ecom_analytics_service/db/mysql"
	"github.com/TranVinhHien/ecom_analytics_service/server"
	"github.com/TranVinhHien/ecom_analytics_service/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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
	conn_order, err := connectDBWithRetry(5, env.DBSourceOrder)
	if conn_order == nil {
		log.Err(err).Msg("Error when created connect to database order")
		return
	}
	// close connection after gin stopped
	defer conn_order.Close()

	log.Info().Msg("Connect to database order successfully")
	db_order := db.NewStoreOrder(conn_order)

	// create connect to database
	conn_transaction, err := connectDBWithRetry(5, env.DBSourceTransaction)
	if conn_transaction == nil {
		log.Err(err).Msg("Error when created connect to database transaction")
		return
	}
	// close connection after gin stopped
	defer conn_transaction.Close()

	log.Info().Msg("Connect to database transaction successfully")
	db_transaction := db.NewStoreTransaction(conn_transaction)

	// create connect to interact database
	conn_interact, err := connectDBWithRetry(5, env.DBSourceInteract)
	if conn_interact == nil {
		log.Err(err).Msg("Error when created connect to database interact")
		return
	}
	// close connection after gin stopped
	defer conn_interact.Close()

	log.Info().Msg("Connect to database interact successfully")
	db_interact := db.NewStoreInteract(conn_interact)

	// create connect to interact database
	conn_agent_ai_db, err := connectDBWithRetry(5, env.DBSourceAgentAIDB)
	if conn_agent_ai_db == nil {
		log.Err(err).Msg("Error when created connect to database interact")
		return
	}
	// close connection after gin stopped
	defer conn_agent_ai_db.Close()

	log.Info().Msg("Connect to database interact successfully")
	db_agent_ai_db := db.NewStoreAgentAIDB(conn_agent_ai_db)

	log.Info().Msg("Creating gin server...")
	// create jwt
	jwtMaker, err := token.NewJWTMaker(env.JWTSecret)
	if err != nil {
		log.Err(err).Msg("Error create JWTMaker")
		return
	}

	// create instance API server
	APIServer := server.NewAPIServices(env, time.Second*10)

	// setup service
	services := services.NewService(db_order, db_transaction, db_interact, db_agent_ai_db, jwtMaker, env, APIServer)
	// setup controller
	controller := controllers.NewAPIController(services, jwtMaker)

	engine := gin.Default()
	engine.MaxMultipartMemory = 32 << 20 // 32 MB

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

	engine.Run(env.HTTPServerAddress)

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
