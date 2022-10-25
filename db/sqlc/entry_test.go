package db_test

import (
	"context"
	db "go-simple-bank/db/sqlc"
	"go-simple-bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) db.Entry {
	account := createRandomAccount(t)

	args := &db.CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomBalance(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), *args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)

	// Check that postgres generates correct values.
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func createEntriesToAccount(t *testing.T, n int) int64 {
	account := createRandomAccount(t)

	args := &db.CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomBalance(),
	}
	for i := 0; i < n; i++ {
		entry, err := testQueries.CreateEntry(context.Background(), *args)
		require.NoError(t, err)
		require.NotEmpty(t, entry)
		require.Equal(t, args.AccountID, entry.AccountID)
		require.Equal(t, args.Amount, entry.Amount)

		// Check that postgres generates correct values.
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)

		args.Amount = util.RandomBalance()
	}

	return args.AccountID
}

func TestEntries(t *testing.T) {

	t.Run("Create entry", func(t *testing.T) {
		createRandomEntry(t)
	})

	t.Run("Get entry", func(t *testing.T) {
		entry := createRandomEntry(t)

		gotEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, gotEntry)
		require.Equal(t, gotEntry.ID, entry.ID)
		require.Equal(t, gotEntry.AccountID, entry.AccountID)
		require.Equal(t, gotEntry.Amount, entry.Amount)
		require.WithinDuration(t, gotEntry.CreatedAt, entry.CreatedAt, time.Second)
	})

	t.Run("List entries by account", func(t *testing.T) {
		accountId := createEntriesToAccount(t, 10)

		args := &db.ListEntriesByAccountParams{
			AccountID: accountId,
			Limit:     5,
			Offset:    5,
		}

		entries, err := testQueries.ListEntriesByAccount(context.Background(), *args)
		require.NoError(t, err)
		require.Len(t, entries, 5)

		for _, entry := range entries {
			require.NotEmpty(t, entry)
		}
	})

}
