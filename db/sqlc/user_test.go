package db_test

import (
	"context"
	db "go-simple-bank/db/sqlc"
	"go-simple-bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) db.User {
	hashedPw, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	args := &db.CreateUserParams{
		Username:       util.RandomOwner(),
		FullName:       util.RandomOwner(),
		HashedPassword: hashedPw,
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), *args)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.HashedPassword, user.HashedPassword)

	// Check that postgres generates correct values.
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}
func TestUsers(t *testing.T) {

	t.Run("Create user", func(t *testing.T) {
		createRandomUser(t)
	})

	t.Run("Get user", func(t *testing.T) {
		user := createRandomUser(t)

		gotUser, err := testQueries.GetUser(context.Background(), user.Username)
		require.NoError(t, err)
		require.NotEmpty(t, gotUser)
		require.Equal(t, gotUser.Username, user.Username)
		require.Equal(t, gotUser.FullName, user.FullName)
		require.Equal(t, gotUser.HashedPassword, user.HashedPassword)
		require.WithinDuration(t, gotUser.CreatedAt, user.CreatedAt, time.Second)
		require.WithinDuration(t, gotUser.PasswordChangedAt, user.PasswordChangedAt, time.Second)
	})

}
