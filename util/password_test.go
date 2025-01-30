package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(10)

	hashedPassword1, err := HashedPassword(password)
	require.NoError(t, err, "no error when hashing password")
	require.NotEmpty(t, hashedPassword1, "hashed password must not be empty")

	err = CheckPassword(hashedPassword1, password)
	require.NoError(t, err, "no error when checking the right password")

	wrongPassword := RandomString(10)
	err = CheckPassword(hashedPassword1, wrongPassword)

	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashedPassword(password)
	require.NoError(t, err, "no error when hashing password")
	require.NotEmpty(t, hashedPassword2, "hashed password must not be empty")

	require.NotEqual(t, hashedPassword1, hashedPassword2, "same password's hash must not be the same")
}
