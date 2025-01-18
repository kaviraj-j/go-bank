package db

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/kaviraj-j/go-bank/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	// Ensure store and accounts are properly created
	require.NotEmpty(t, store, "Store should not be empty; database transaction handler must be initialized")
	require.NotEmpty(t, fromAccount, "FromAccount should not be empty; a valid account must be created")
	require.NotEmpty(t, toAccount, "ToAccount should not be empty; a valid account must be created")

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
		require.NoError(t, err, "TransferTx should execute without errors")
		require.NotEmpty(t, result, "TransferTx result should not be empty; transaction must return valid data")

		// Check transfers
		transfer := result.Transfer
		require.NotEmpty(t, transfer, "Transfer record should not be empty; a valid transfer must be created")
		require.Equal(t, fromAccount.ID, transfer.FromAccountID, "Transfer should originate from the correct FromAccountID")
		require.Equal(t, toAccount.ID, transfer.ToAccountID, "Transfer should be sent to the correct ToAccountID")
		require.NotZero(t, transfer.ID, "Transfer ID should not be zero; must be stored in the database")
		require.NotZero(t, transfer.CreatedAt, "Transfer should have a valid timestamp")

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err, "Transfer should be retrievable from the database")

		// Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry, "FromEntry should not be empty; a debit entry must be recorded")
		require.Equal(t, fromAccount.ID, fromEntry.AccountID, "FromEntry should belong to the correct FromAccountID")
		require.NotZero(t, fromEntry.ID, "FromEntry ID should not be zero; must be stored in the database")
		require.NotZero(t, fromEntry.CreatedAt, "FromEntry should have a valid timestamp")

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err, "FromEntry should be retrievable from the database")

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry, "ToEntry should not be empty; a credit entry must be recorded")
		require.Equal(t, toAccount.ID, toEntry.AccountID, "ToEntry should belong to the correct ToAccountID")
		require.NotZero(t, toEntry.ID, "ToEntry ID should not be zero; must be stored in the database")
		require.NotZero(t, toEntry.CreatedAt, "ToEntry should have a valid timestamp")

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err, "ToEntry should be retrievable from the database")

		// Check account balances
		transferedFromAccount := result.FromAccount
		require.NotEmpty(t, transferedFromAccount, "TransferredFromAccount should not be empty; account data should be returned after transfer")
		require.Equal(t, fromAccount.ID, transferedFromAccount.ID, "TransferredFromAccount ID should match original FromAccount ID")

		transferedToAccount := result.ToAccount
		require.NotEmpty(t, transferedToAccount, "TransferredToAccount should not be empty; account data should be returned after transfer")
		require.Equal(t, toAccount.ID, transferedToAccount.ID, "TransferredToAccount ID should match original ToAccount ID")

		// Convert balances from string to decimal.Decimal
		fromAccountBalance, err := decimal.NewFromString(fromAccount.Balance)
		require.NoError(t, err, "FromAccount balance should be a valid decimal")

		transferedFromAccountBalance, err := decimal.NewFromString(transferedFromAccount.Balance)
		require.NoError(t, err, "TransferredFromAccount balance should be a valid decimal")

		toAccountBalance, err := decimal.NewFromString(toAccount.Balance)
		require.NoError(t, err, "ToAccount balance should be a valid decimal")

		transferedToAccountBalance, err := decimal.NewFromString(transferedToAccount.Balance)
		require.NoError(t, err, "TransferredToAccount balance should be a valid decimal")

		// Ensure the balance changes match
		diff1 := fromAccountBalance.Sub(transferedFromAccountBalance).String()
		diff2 := transferedToAccountBalance.Sub(toAccountBalance).String()

		fmt.Printf("=================\nfromAccountBalance: %v\ntransferedFromAccountBalance: %v\ntoAccountBalance: %v\ntransferedToAccountBalance: %v\ndiff1: %v\ndiff2: %v\n=======================\n", fromAccountBalance, transferedFromAccountBalance, toAccountBalance, transferedToAccountBalance, diff1, diff2)

		require.Equal(t, diff1, diff2, "The amount deducted from FromAccount should be equal to the amount added to ToAccount")

		// Uncomment if strict balance checking is required
		// amount, err := decimal.NewFromString(transfer.Amount)
		// require.NoError(t, err, "Transfer amount should be a valid decimal")
		// expectedFromBalance := fromAccountBalance.Sub(amount)
		// expectedToBalance := toAccountBalance.Add(amount)
		// require.Equal(t, expectedFromBalance, transferedFromAccountBalance, "FromAccount balance mismatch")
		// require.Equal(t, expectedToBalance, transferedToAccountBalance, "ToAccount balance mismatch")
	}
}
