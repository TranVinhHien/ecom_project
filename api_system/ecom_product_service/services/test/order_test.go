package services_test

// import (
// 	"context"
// 	"testing"

// 	services "github.com/TranVinhHien/ecom_product_service/services/entity"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/stretchr/testify/require"
// )

// func TestCreateOrder(t *testing.T) {
// 	// Create a test context
// 	ctx := context.Background()
// 	_, err := testService.CreateOrder(ctx, "9930", &services.CreateOrderParams{
// 		Address_id:  "ads-104933",
// 		Discount_Id: "1030f88e-dfb2-44cb-85af-38f06e7c65d5",
// 		Payment_id:  "550e8400-e29b-41d4-a716-446655440051",
// 		NumOfProducts: []services.AmountProdduct{
// 			{Product_sku_id: "100004", Amount: 2}, //6721
// 			{Product_sku_id: "100005", Amount: 1}, //1738
// 			{Product_sku_id: "100006", Amount: 1}, //3591
// 		},
// 	})
// 	// require.NotEqual(t, err.Err.Error(), "")

// 	// fmt.Println("categories:", categories)

// 	// require.Empty(t, categories)

// 	require.Nil(t, err.Err)

// 	require.True(t, true, "Test category not found in retrieved categories")
// }
// func TestSubmitOrder(t *testing.T) {
// 	// Create a test context
// 	ctx := context.Background()
// 	testService.CallBackMoMo(ctx, services.TransactionMoMO{
// 		OrderID:    "ec9f2a06-3331-49d7-951b-4b9424679914",
// 		ResultCode: 0,
// 	})
// 	// require.NotEqual(t, err.Err.Error(), "")

// 	// fmt.Println("categories:", categories)

// 	// require.Empty(t, categories)

// 	require.NotNil(t, nil)

// 	require.True(t, true, "Test category not found in retrieved categories")
// }

// func TestListOrderByUserID(t *testing.T) {
// 	// Create a test context
// 	ctx := context.Background()
// 	user_id := "9930"
// 	orders, err := testService.ListOrderByUserID(ctx, user_id, services.QueryFilter{
// 		Conditions: []services.Condition{
// 			{Field: "total_amount", Operator: ">=", Value: 10},
// 			// {Field: "end_date", Operator: ">=", Value: time.Now()},
// 			// {Field: "start_date", Operator: "=", Value: time.Now()},
// 		},
// 		// OrderBy: &services.OrderBy{
// 		// 	// Field: "create_date",
// 		// 	// Value: "ASC",
// 		// },
// 		Page:     1,
// 		PageSize: 10,
// 	})
// 	require.Nil(t, err)
// 	require.Empty(t, orders)

// 	require.True(t, true, "Test category not found in retrieved orders")
// }
