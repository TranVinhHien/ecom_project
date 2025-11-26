package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration) (*Payload, string, error)
	VerifyToken(token string) (*Payload, error)
}
