package services

import (
	config_assets "github.com/TranVinhHien/ecom_analytics_service/assets/config"
	"github.com/TranVinhHien/ecom_analytics_service/assets/token"
	db_interact "github.com/TranVinhHien/ecom_analytics_service/db/mysql"
	db_order "github.com/TranVinhHien/ecom_analytics_service/db/mysql"
	db_transaction "github.com/TranVinhHien/ecom_analytics_service/db/mysql"
	"github.com/TranVinhHien/ecom_analytics_service/server"
)

type service struct {
	transaction    db_transaction.StoreTransaction
	order          db_order.StoreOrder
	interact       db_interact.StoreInteract
	db_agent_ai_db db_interact.StoreAgentAIDB
	jwt            token.Maker
	env            config_assets.ReadENV
	apiServer      server.ApiServer
}

func NewService(order db_order.StoreOrder, transaction db_transaction.StoreTransaction, interact db_interact.StoreInteract, db_agent_ai_db db_interact.StoreAgentAIDB, jwt token.Maker, env config_assets.ReadENV, apiServer server.ApiServer) ServiceUseCase {
	return &service{order: order, transaction: transaction, interact: interact, db_agent_ai_db: db_agent_ai_db, jwt: jwt, env: env, apiServer: apiServer}
}
