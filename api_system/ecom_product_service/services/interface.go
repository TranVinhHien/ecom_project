package services

import (
	"context"

	services "github.com/TranVinhHien/ecom_product_service/services/entity"
	iservices "github.com/TranVinhHien/ecom_product_service/services/interface"
)

//	type ServicesRepository interface {
//		iservices.CategoriesRepository
//	}
type ServiceUseCase interface {
	iservices.Categories
	iservices.Products
	iservices.Media
}

type ServicesRedis interface {
	StartExpirationListenerOrderOnline(func(ctx context.Context, orderID string))

	CheckExistsFromBlackList(ctx context.Context, token string, exprid float64) bool
	RemoveTokenExp(zsetKey string)
	AddTokenToBlackList(ctx context.Context, token string, exprid float64) error
	// category
	AddCategories(ctx context.Context, cates []services.Categorys) error
	RemoveCategories(ctx context.Context) error
	GetCategoryTree(ctx context.Context, rootID string) ([]services.Categorys, error)

	DeleteOrderOnline(ctx context.Context, orderID string) error
}
