package db_test

import (
	"context"
	"database/sql"
	db "go-simple-bank/db/sqlc"
	"go-simple-bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) db.Account {
	user := createRandomUser(t)

	args := &db.CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), *args)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	// Check that postgres generates correct values.
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
func TestAccounts(t *testing.T) {

	t.Run("Create account", func(t *testing.T) {
		createRandomAccount(t)
	})

	t.Run("Get Account", func(t *testing.T) {
		account := createRandomAccount(t)

		gotAccount, err := testQueries.GetAccount(context.Background(), account.ID)
		require.NoError(t, err)
		require.NotEmpty(t, gotAccount)
		require.Equal(t, gotAccount.Owner, account.Owner)
		require.Equal(t, gotAccount.Balance, account.Balance)
		require.Equal(t, gotAccount.Currency, account.Currency)
		require.WithinDuration(t, account.CreatedAt, gotAccount.CreatedAt, time.Second)
	})

	t.Run("List Accounts", func(t *testing.T) {
		var lastAccount db.Account
		for i := 0; i < 10; i++ {
			lastAccount = createRandomAccount(t)
		}

		args := &db.ListAccountsParams{
			Owner:  lastAccount.Owner,
			Limit:  5,
			Offset: 0,
		}

		accounts, err := testQueries.ListAccounts(context.Background(), *args)
		require.NoError(t, err)
		require.NotEmpty(t, accounts)

		for _, account := range accounts {
			require.NotEmpty(t, account)
			require.Equal(t, lastAccount.Owner, account.Owner)
		}
	})

	t.Run("Update Account", func(t *testing.T) {
		account := createRandomAccount(t)

		args := &db.UpdateAccountParams{
			ID:      account.ID,
			Balance: account.Balance * 2,
		}
		err := testQueries.UpdateAccount(context.Background(), *args)
		require.NoError(t, err)

		updatedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
		require.NotEmpty(t, updatedAccount)
		require.Equal(t, updatedAccount.Owner, account.Owner)
		require.Equal(t, updatedAccount.Balance, account.Balance*2)
		require.Equal(t, updatedAccount.Currency, account.Currency)
		require.WithinDuration(t, account.CreatedAt, updatedAccount.CreatedAt, time.Second)
	})

	t.Run("Delete Account", func(t *testing.T) {
		account := createRandomAccount(t)

		err := testQueries.DeleteAccount(context.Background(), account.ID)
		require.NoError(t, err)

		gotAccount, err := testQueries.GetAccount(context.Background(), account.ID)
		require.Error(t, err)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Empty(t, gotAccount)
	})
}
