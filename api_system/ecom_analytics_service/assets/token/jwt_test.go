package token

// write a function test function CreateToken and VerifyToken

// func TestCreateToken(t *testing.T) {
// 	maker, err := NewJWTMaker("secret-key")
// 	require.Nil(t, err)

// 	duration, _ := time.ParseDuration("1h")
// 	payload, token, err := maker.CreateToken("john.doe", duration)
// 	log.Printf("payload: %v", payload)
// 	log.Printf("valid: %v", payload.Valid())
// 	log.Printf("token: %s", token)
// 	require.Nil(t, err)

// 	require.NotNil(t, payload)
// 	require.NotEmpty(t, token)
// 	// test verification
// 	verifiedPayload, err := maker.VerifyToken(token)
// 	require.Nil(t, err)
// 	require.Equal(t, payload.Sub, verifiedPayload.Sub)
// 	require.Equal(t, payload.Exp, verifiedPayload.Exp)
// }
