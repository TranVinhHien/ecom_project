package services_test

// import (
// 	"context"
// 	"database/sql"
// 	"os"
// 	"testing"

// 	config_assets "github.com/TranVinhHien/ecom_product_service/assets/config"
// 	assets_firebase "github.com/TranVinhHien/ecom_product_service/assets/fire-base"
// 	assets_jobs "github.com/TranVinhHien/ecom_product_service/assets/jobs"
// 	"github.com/TranVinhHien/ecom_product_service/assets/token"
// 	db "github.com/TranVinhHien/ecom_product_service/db/mysql"
// 	redis_db "github.com/TranVinhHien/ecom_product_service/db/redis"
// 	"github.com/TranVinhHien/ecom_product_service/services"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/redis/go-redis/v9"
// )

// var testService services.ServiceUseCase

// func TestMain(m *testing.M) {

// 	pool, _ := sql.Open("mysql", "root:12345@tcp(localhost:3306)/e-commerce?parseTime=true")

// 	env, _ := config_assets.LoadConfig("../../")
// 	db := db.NewStore(pool)
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// 	rerids := redis_db.NewRedisDB(rdb)
// 	jwtMaker, _ := token.NewJWTMaker(env.JWTSecret)
// 	// firebaseApp, _ := firebase.N(context.Background(), nil, )
// 	firebase, _ := assets_firebase.NewFirebase(context.Background(), "../../assets/firebase/credentials.json")
// 	job, _ := assets_jobs.NewJobScheduler()
// 	testService = services.NewService(db, jwtMaker, env, rerids, firebase, job)

// 	os.Exit(m.Run())
// }
