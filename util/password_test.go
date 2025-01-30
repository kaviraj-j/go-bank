package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := randomString(10)
	hashedPassword, err := HashedPassword(password)

	require.NoError(t, err, "no error when hashing password")
	require.NotEmpty(t, hashedPassword, "hashed password must not be empty")

	err = CheckPassword(hashedPassword, password)
	require.NoError(t, err, "no error when checking the right password")

	wrongPassword := randomString(10)
	err = CheckPassword(hashedPassword, wrongPassword)

	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

}
