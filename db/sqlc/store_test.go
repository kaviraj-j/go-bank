package db

import (
	"context"
	"sync"
	"testing"

	"github.com/kaviraj-j/go-bank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	require.NotEmpty(t, store)
	require.NotEmpty(t, fromAccount)
	require.NotEmpty(t, toAccount)

	n := 5

	type ResultChannel struct {
		result TransferTxResult
		err    error
	}

	results := make(chan ResultChannel, n)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        util.RandomAmount(),
			})
			results <- ResultChannel{
				result: result,
				err:    err,
			}
		}()
	}

	// Close the results channel after all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Read results from the channel
	for res := range results {
		err := res.err
		result := res.result
		require.NoError(t, err)

		require.NotEmpty(t, result)

		// Check transfers
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotZero(t, fromEntry.AccountID)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.AccountID)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
	}
}
