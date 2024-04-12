package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	var (
		ctx             = context.Background()
		store           = NewStore(testDB)
		transferAccount = createRandomAccount(t)
		receiveAccount  = createRandomAccount(t)
		// run n concurrent transfer transactions
		n      = 5
		amount = int64(10)

		errs    = make(chan error)
		results = make(chan TransferTxResult)
	)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: transferAccount.ID,
				ToAccountID:   receiveAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transferAccount.ID, transfer.FromAccountID)
		require.Equal(t, receiveAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(ctx, transfer.ID)
		require.NoError(t, err)

		//check entries
		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, transferAccount.ID, result.FromEntry.AccountID)
		require.Equal(t, receiveAccount.ID, result.ToEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)
		require.NotZero(t, result.ToEntry.CreatedAt)

		_, err = store.GetEntry(ctx, result.FromEntry.ID)
		require.NoError(t, err)

		_, err = store.GetEntry(ctx, result.ToEntry.ID)
		require.NoError(t, err)

		//TODO: check account's balance

	}

}
