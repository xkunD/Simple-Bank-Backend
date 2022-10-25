package db_test

import (
	"context"
	db "go-simple-bank/db/sqlc"
	"go-simple-bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) db.Transfer {
	account := createRandomAccount(t)
	account2 := createRandomAccount(t)

	args := &db.CreateTransferParams{
		FromAccountID: account.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomBalance(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), *args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.Amount, transfer.Amount)

	// Check that postgres generates correct values.
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func createTransfersBetweenAccounts(t *testing.T, n int) (int64, int64) {
	account := createRandomAccount(t)
	account2 := createRandomAccount(t)

	args := &db.CreateTransferParams{
		FromAccountID: account.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomBalance(),
	}

	for i := 0; i < n; i++ {
		transfer, err := testQueries.CreateTransfer(context.Background(), *args)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)
		require.Equal(t, args.FromAccountID, transfer.FromAccountID)
		require.Equal(t, args.ToAccountID, transfer.ToAccountID)
		require.Equal(t, args.Amount, transfer.Amount)

		// Check that postgres generates correct values.
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		args.Amount = util.RandomBalance()
	}

	return args.FromAccountID, args.ToAccountID
}
func TestTransfers(t *testing.T) {

	t.Run("Create transfer", func(t *testing.T) {
		createRandomTransfer(t)
	})

	t.Run("Get transfer", func(t *testing.T) {
		transfer := createRandomTransfer(t)

		gotTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		require.NotEmpty(t, gotTransfer)
		require.Equal(t, gotTransfer.ID, transfer.ID)
		require.Equal(t, gotTransfer.FromAccountID, transfer.FromAccountID)
		require.Equal(t, gotTransfer.ToAccountID, transfer.ToAccountID)
		require.Equal(t, gotTransfer.Amount, transfer.Amount)
		require.WithinDuration(t, gotTransfer.CreatedAt, transfer.CreatedAt, time.Second)
	})

	t.Run("List transfers by accounts", func(t *testing.T) {
		fromAccountId, toAccountId := createTransfersBetweenAccounts(t, 10)

		args := &db.ListTransfersByAccountParams{
			FromAccountID: fromAccountId,
			ToAccountID:   toAccountId,
			Limit:         5,
			Offset:        5,
		}

		transfers, err := testQueries.ListTransfersByAccount(context.Background(), *args)
		require.NoError(t, err)
		require.Len(t, transfers, 5)

		for _, transfer := range transfers {
			require.NotEmpty(t, transfer)
		}
	})

}
