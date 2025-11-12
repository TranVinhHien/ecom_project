package iservices

import (
	"context"

	assets_services "github.com/TranVinhHien/ecom_payment-service/services/assets"
	services "github.com/TranVinhHien/ecom_payment-service/services/entity"
)

type Payments interface {
	PaymentMethodDetail(ctx context.Context, id string) (map[string]interface{}, *assets_services.ServiceError)
	ListPaymentMethod(ctx context.Context) (map[string]interface{}, *assets_services.ServiceError)
	InitPayment(ctx context.Context, userId string, email string, order services.InitPaymentParams) (map[string]interface{}, *assets_services.ServiceError)
	CallBackMoMo(ctx context.Context, tran services.TransactionMoMO)
	GetURLOrderMoMOAgain(ctx context.Context, user_id string) (map[string]interface{}, *assets_services.ServiceError)
	// (ctx context.Context, user_id string) (map[string]interface{}, *assets_services.ServiceError)
	// GetURLOrderMoMOAgain(ctx context.Context, user_id string) (map[string]interface{}, *assets_services.ServiceError)
}
type Jobs interface {
	CheckTransactionTimeout(ctx context.Context)
}
