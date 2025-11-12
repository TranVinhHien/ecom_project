package services_test

// import (
// 	"context"
// 	"testing"

// 	services "github.com/TranVinhHien/ecom_analytics_service/services/entity"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/stretchr/testify/require"
// )

// func TestGetDiscount(t *testing.T) {
// 	// Create a test context
// 	ctx := context.Background()

// 	categories, _ := testService.ListDiscount(ctx, services.QueryFilter{
// 		Conditions: []services.Condition{
// 			// {Field: "amount", Operator: ">=", Value: 1},
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
// 	// require.NotEqual(t, err.Err.Error(), "")

// 	// fmt.Println("categories:", categories)

// 	require.Empty(t, categories)

// 	require.True(t, true, "Test category not found in retrieved categories")
// }
