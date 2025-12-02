package token

import (
	"errors"
	"time"
)

var (
	ErrExpireToken  = errors.New("token has Expired")
	ErrInvalidToken = errors.New("token has Invalid")
)

type Payload struct {
	Sub    string `json:"sub"`
	Iss    string `json:"iss"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
	Scope  string `json:"scope"`
	Jti    string `json:"jti"`
	UserId string `json:"userId"`
	Email  string `json:"email"`
}

func CreateNewPayload(username string, duration time.Duration) *Payload {
	payload := Payload{
		Sub: username,
		Iss: "hienlazada.edu.vn",
		// Aud: "sau này truyền role vào đây",
		Exp: time.Now().Add(duration).Unix(),
		Iat: time.Now().Unix(),
	}
	return &payload
}

func (payload *Payload) Valid() bool {
	currentTime := time.Now().Unix()
	return currentTime <= payload.Exp
}
