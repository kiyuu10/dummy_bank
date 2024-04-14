package db

import (
	"context"
	"fmt"
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
		n      = 2
		amount = int64(10)

		errs    = make(chan error)
		results = make(chan TransferTxResult)
	)

	fmt.Println(">> before >> transfer account: ", transferAccount.Balance, " - Receive account: ", receiveAccount.Balance)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d: ", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
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
	exitsted := make(map[int]bool)
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

		// check accounts
		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, transferAccount.ID, result.FromAccount.ID)

		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, receiveAccount.ID, result.ToAccount.ID)

		//TODO: check account's balance
		fmt.Println(">> tx >> transfer account: ", result.FromAccount.Balance, " - Receive account: ", result.ToAccount.Balance)
		balanceTransferAcc := transferAccount.Balance - result.FromAccount.Balance
		balanceReceiveAcc := result.ToAccount.Balance - receiveAccount.Balance
		require.Equal(t, balanceTransferAcc, balanceReceiveAcc)
		require.True(t, balanceTransferAcc > 0)
		require.True(t, balanceTransferAcc%amount == 0)

		k := int(balanceTransferAcc / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, exitsted, k)
		exitsted[k] = true

		// check the final updated balances
		accTransferUpdated, err := testQueries.GetAccount(ctx, transferAccount.ID)
		require.NoError(t, err)

		accReceiveUpdated, err := testQueries.GetAccount(ctx, receiveAccount.ID)
		require.NoError(t, err)

		fmt.Println(">> after updated>> transfer account: ", accTransferUpdated.Balance, " - Receive account: ", accReceiveUpdated.Balance)
		fmt.Println(">> after >> transfer account: ", transferAccount.Balance, " - Receive account: ", receiveAccount.Balance)
		require.Equal(t, transferAccount.Balance-(int64(n)*amount), accTransferUpdated.Balance)
		require.Equal(t, receiveAccount.Balance+(int64(n)*amount), accReceiveUpdated.Balance)

	}

}
