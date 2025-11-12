package services_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	services "github.com/TranVinhHien/ecom_analytics_service/services/entity"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/stretchr/testify/require"
// )

// func TestGetALL(t *testing.T) {
// 	// Create a test context
// 	ctx := context.Background()

// 	categories, _ := testService.GetAllProductSimple(ctx, services.QueryFilter{
// 		Conditions: []services.Condition{
// 			{Field: "category_id", Operator: "=", Value: 10},
// 			// {Field: "end_date", Operator: ">=", Value: time.Now()},
// 			// {Field: "start_date", Operator: "=", Value: time.Now()},
// 		},
// 		OrderBy: &services.OrderBy{
// 			Field: "name",
// 			Value: "ASC",
// 		},
// 		Page:     1,
// 		PageSize: 2,
// 	})
// 	// require.NotEqual(t, err.Err.Error(), "")

// 	fmt.Println("categories:", categories)

// 	require.Empty(t, categories)

// 	require.True(t, true, "Test category not found in retrieved categories")
// }
