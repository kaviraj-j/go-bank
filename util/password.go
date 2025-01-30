package util

import "golang.org/x/crypto/bcrypt"

// returns the bcrypt hash of a raw password
func HashedPassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// checks if the provided password is correct or not
func CheckPassword(hashedPassword string, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
}
