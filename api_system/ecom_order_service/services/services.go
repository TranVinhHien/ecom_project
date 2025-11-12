package services

import (
	config_assets "github.com/TranVinhHien/ecom_order_service/assets/config"
	"github.com/TranVinhHien/ecom_order_service/assets/token"
	db "github.com/TranVinhHien/ecom_order_service/db/mysql"
	"github.com/TranVinhHien/ecom_order_service/server"
)

type service struct {
	repository db.Store
	redis      ServicesRedis
	jwt        token.Maker
	env        config_assets.ReadENV
	apiServer  server.ApiServer
	// firebase   *assets_firebase.FirebaseMessaging
	// jobs       *assets_jobs.JobScheduler
}

func NewService(repo db.Store, jwt token.Maker, env config_assets.ReadENV, redis ServicesRedis, apiServer server.ApiServer) ServiceUseCase {
	return &service{repository: repo, jwt: jwt, env: env, redis: redis, apiServer: apiServer}
}
