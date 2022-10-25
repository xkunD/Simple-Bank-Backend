package db_test

import (
	"context"
	db "go-simple-bank/db/sqlc"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {

	t.Run("Transfer TX", func(t *testing.T) {
		store := db.NewStore(testDb)

		accountFrom := createRandomAccount(t)
		accountTo := createRandomAccount(t)
		log.Printf(">> BEFORE: FROM(%d), TO(%d) ", accountFrom.Balance, accountTo.Balance)

		// run n concurrent transfer transactions
		n := 5
		amount := int64(10)

		// Channel to communicate errors to testing function becuase transfer is inside go routine.
		errs := make(chan error)
		// Channel to receive results.
		results := make(chan db.TransferTxResult)

		for i := 0; i < n; i++ {
			go func() {
				result, err := store.TransferTx(context.Background(), db.TransferTxParams{
					FromAccountId: accountFrom.ID,
					ToAccountId:   accountTo.ID,
					Amount:        amount,
				})
				errs <- err
				results <- result
			}()
		}

		existed := make(map[int]bool)
		// checks from outside go routine.
		for i := 0; i < n; i++ {
			err := <-errs
			require.NoError(t, err)

			// checks transfer
			results := <-results
			require.NotEmpty(t, results)
			transfer := results.Transfer
			require.Equal(t, transfer.FromAccountID, accountFrom.ID)
			require.Equal(t, transfer.ToAccountID, accountTo.ID)
			require.Equal(t, transfer.Amount, amount)
			require.NotZero(t, transfer.ID)
			require.NotZero(t, transfer.CreatedAt)
			_, err = store.GetTransfer(context.Background(), transfer.ID)
			require.NoError(t, err)

			// checks entries
			require.NotEmpty(t, results.FromEntry)
			fromEntry := results.FromEntry
			require.Equal(t, fromEntry.AccountID, accountFrom.ID)
			require.Equal(t, fromEntry.Amount, -amount)
			require.NotZero(t, fromEntry.ID)
			require.NotZero(t, fromEntry.CreatedAt)
			_, err = store.GetEntry(context.TODO(), fromEntry.ID)
			require.NoError(t, err)

			require.NotEmpty(t, results.ToEntry)
			toEntry := results.ToEntry
			require.Equal(t, toEntry.AccountID, accountTo.ID)
			require.Equal(t, toEntry.Amount, amount)
			require.NotZero(t, toEntry.ID)
			require.NotZero(t, toEntry.CreatedAt)
			_, err = store.GetEntry(context.TODO(), toEntry.ID)
			require.NoError(t, err)

			// check accounts
			require.NotEmpty(t, results.FromAccount)
			fromAccount := results.FromAccount
			require.Equal(t, fromAccount.ID, accountFrom.ID)

			require.NotEmpty(t, results.ToAccount)
			toAccount := results.ToAccount
			require.Equal(t, toAccount.ID, accountTo.ID)

			//check balances
			log.Printf(">> TX: FROM(%d), TO(%d) ", fromAccount.Balance, toAccount.Balance)
			diffFromAccount := accountFrom.Balance - fromAccount.Balance
			diffToAccount := toAccount.Balance - accountTo.Balance
			require.Equal(t, diffFromAccount, diffToAccount)
			require.True(t, diffFromAccount > 0)
			require.True(t, diffToAccount > 0)

			require.True(t, diffFromAccount%amount == 0) // 1 * amount, 2 * amount, ... , n * amount

			k := int(diffFromAccount / amount)
			require.True(t, k >= 1 && k <= n)
			require.NotContains(t, existed, k)
			existed[k] = true
		}

		//check final updated balances
		updatedFromAccount, err := testQueries.GetAccount(context.Background(), accountFrom.ID)
		require.NoError(t, err)
		require.NotEmpty(t, updatedFromAccount)

		updatedToAccount, err := testQueries.GetAccount(context.Background(), accountTo.ID)
		require.NoError(t, err)
		require.NotEmpty(t, updatedToAccount)

		log.Printf(">> AFTER: FROM(%d), TO(%d) ", updatedFromAccount.Balance, updatedToAccount.Balance)
		require.Equal(t, accountFrom.Balance-int64(n)*amount, updatedFromAccount.Balance)
		require.Equal(t, accountTo.Balance+int64(n)*amount, updatedToAccount.Balance)
	})

	t.Run("Transfer TX Deadlock", func(t *testing.T) {
		store := db.NewStore(testDb)

		accountFrom := createRandomAccount(t)
		accountTo := createRandomAccount(t)
		log.Printf(">> BEFORE: FROM(%d), TO(%d) ", accountFrom.Balance, accountTo.Balance)

		// run n concurrent transfer transactions
		n := 10
		amount := int64(10)

		// Channel to communicate errors to testing function becuase transfer is inside go routine.
		errs := make(chan error)

		for i := 0; i < n; i++ {
			fromAccount := accountFrom.ID
			toAccount := accountTo.ID
			if i%2 == 0 {
				fromAccount = accountTo.ID
				toAccount = accountFrom.ID
			}

			go func() {
				_, err := store.TransferTx(context.Background(), db.TransferTxParams{
					FromAccountId: fromAccount,
					ToAccountId:   toAccount,
					Amount:        amount,
				})
				errs <- err
			}()
		}

		for i := 0; i < n; i++ {
			err := <-errs
			require.NoError(t, err)

		}

		//check final updated balances
		updatedFromAccount, err := testQueries.GetAccount(context.Background(), accountFrom.ID)
		require.NoError(t, err)
		require.NotEmpty(t, updatedFromAccount)

		updatedToAccount, err := testQueries.GetAccount(context.Background(), accountTo.ID)
		require.NoError(t, err)
		require.NotEmpty(t, updatedToAccount)

		log.Printf(">> AFTER: FROM(%d), TO(%d) ", updatedFromAccount.Balance, updatedToAccount.Balance)
		require.Equal(t, accountFrom.Balance, updatedFromAccount.Balance)
		require.Equal(t, accountTo.Balance, updatedToAccount.Balance)

	})

}
