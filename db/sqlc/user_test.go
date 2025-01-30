package db

import (
	"context"
	"testing"

	"github.com/kaviraj-j/go-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	randomPassword := util.RandomString(10)
	hashedPassword, err := util.HashedPassword(randomPassword)
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.PasswordChangedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	newUser := createRandomUser(t)
	user1, err := testQueries.GetUser(context.Background(), newUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user1)

	require.Equal(t, newUser.Username, user1.Username)
	require.Equal(t, newUser.FullName, user1.FullName)
	require.Equal(t, newUser.Email, user1.Email)
	require.Equal(t, newUser.HashedPassword, user1.HashedPassword)

	require.NotZero(t, user1.CreatedAt)
	require.NotZero(t, user1.PasswordChangedAt)
}
