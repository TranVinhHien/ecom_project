package db_test

// import (
// 	"context"
// 	"database/sql"
// 	db "github.com/TranVinhHien/ecom_payment_service/db/sqlc"
// 	"testing"
// 	"time"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/stretchr/testify/require"
// )

// func createRandomUser() db.Users {
// 	return db.Users{
// 		Username: "testuser",
// 		Password: "password",
// 		FullName: "Test User",
// 		IsActive: true,
// 		CreateAt: time.Now(),
// 	}
// }

// func TestGetUser(t *testing.T) {
// 	pool, err := sql.Open("mysql", "root:12345@tcp(localhost:3306)/e-commerce")
// 	store := db.New(pool)

// 	// Create a random user
// 	user := createRandomUser()

// 	// Insert the user into the database

// 	require.NoError(t, err)

// 	// Retrieve the user from the database
// 	retrievedUser, err := store.GetUser(context.Background(), user.Username)
// 	require.NoError(t, err)

// 	// Verify the retrieved user's information
// 	require.Equal(t, user.Username, retrievedUser.Username)
// 	require.Equal(t, user.Password, retrievedUser.Password)
// 	require.Equal(t, user.FullName, retrievedUser.FullName)
// 	require.Equal(t, user.IsActive, retrievedUser.IsActive)
// 	require.WithinDuration(t, user.CreateAt, retrievedUser.CreateAt, time.Second)
// }
