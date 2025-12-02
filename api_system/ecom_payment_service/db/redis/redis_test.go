package redis_db

// write function test removeExpiredTokens and addToken

// func TestRedisDBaddScoreMember(t *testing.T) {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// 	db := NewRedisDB(rdb)
// 	ctx := context.Background()
// 	err := db.addScoreMember(ctx, "token_set", "test_token1", float64(time.Now().Add(time.Second*15).Unix()))
// 	require.NoError(t, err)
// 	err = db.addScoreMember(ctx, "token_set", "test_token2", float64(time.Now().Add(time.Second*15).Unix()))
// 	require.NoError(t, err)
// }
// func TestRedisDBisExists(t *testing.T) {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// 	db := NewRedisDB(rdb)
// 	ctx := context.Background()
// 	check := db.isExists(ctx, "token_set", "test_token1")
// 	fmt.Println(check, "token1")
// 	check = db.isExists(ctx, "token_set", "test_token2")
// 	fmt.Println(check, "token2")
// 	require.NotNil(t, nil)
// }
// func TestRedisDBremove(t *testing.T) {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// 	db := NewRedisDB(rdb)
// 	ctx := context.Background()
// 	err := db.removeExpired(ctx, "token_set")
// 	require.NoError(t, err)
// }
// func TestRedisStart(t *testing.T) {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// 	db := NewRedisDB(rdb)
// 	db.RemoveTokenExp("token_set")
// }
