package services

import (
	"context"
	"time"

	db "github.com/TranVinhHien/ecom_payment-service/db/sqlc"
	services "github.com/TranVinhHien/ecom_payment-service/services/entity"
	iservices "github.com/TranVinhHien/ecom_payment-service/services/interface"
)

type ServicesRepository interface {
	// iservices.UserRepository
	// iservices.ProductRepository
	// iservices.CategoriesRepository
	// iservices.ĐiscountRepository
	// iservices.PaymentRepository
	// iservices.OrderRepository
	// iservices.RatingRepository
	// iservices.ProductsRepository
	db.Querier
}
type ServiceUseCase interface {
	// iservices.Order
	iservices.Payments
	iservices.Jobs
}

type ServicesRedis interface {

	// orderOnline
	AddTransactionOnline(ctx context.Context, user_id string, payload services.CombinedDataPayLoadMoMo, duration time.Duration) error
	GetTransactionOnline(ctx context.Context, user_id string) (payload *services.CombinedDataPayLoadMoMo, err error)
	DeleteTransactionOnline(ctx context.Context, transactionID string) error
	GetTransactionOnlineWithIDTran(ctx context.Context, transactionID string) (payload *services.CombinedDataPayLoadMoMo, err error)
}
