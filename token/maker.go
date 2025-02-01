package token

import "time"

type Maker interface {

	// CreateToken creates token with a valid duration and username
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
