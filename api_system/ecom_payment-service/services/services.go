package services

import (
	config_assets "github.com/TranVinhHien/ecom_payment-service/assets/config"
	"github.com/TranVinhHien/ecom_payment-service/assets/email"
	"github.com/TranVinhHien/ecom_payment-service/assets/token"
	db "github.com/TranVinhHien/ecom_payment-service/db/mysql"
	"github.com/TranVinhHien/ecom_payment-service/kafka"
	"github.com/TranVinhHien/ecom_payment-service/server"
)

type service struct {
	redis ServicesRedis
	jwt   token.Maker
	env   config_assets.ReadENV
	// jobs  *assets_jobs.JobScheduler

	repository db.Store
	apiServer  server.ApiServer
	producer   kafka.EventProducer
	email      email.BrevoEmailService
}

func NewService(jwt token.Maker, env config_assets.ReadENV, redis ServicesRedis, repository db.Store, apiServer server.ApiServer, producer kafka.EventProducer) ServiceUseCase {
	// Khởi tạo event publisher từ kafka client
	// eventPublisher := kafka.(producer)

	return &service{
		jwt:        jwt,
		env:        env,
		redis:      redis,
		repository: repository,
		apiServer:  apiServer,
		producer:   producer,
		email:      *email.NewBrevoEmailService(env.BrevoAPIKey, env.SenderEmail, env.SenderName),
	}
}
