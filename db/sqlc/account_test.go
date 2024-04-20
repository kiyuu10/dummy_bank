package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kiyuu10/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	inputAccount := createRandomAccount(t)
	expectAccount, err := testQueries.GetAccount(context.Background(), inputAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, expectAccount)

	require.Equal(t, inputAccount.ID, expectAccount.ID)
	require.Equal(t, inputAccount.Owner, expectAccount.Owner)
	require.Equal(t, inputAccount.Balance, expectAccount.Balance)
	require.Equal(t, inputAccount.Currency, expectAccount.Currency)
	require.WithinDuration(t, inputAccount.CreatedAt, expectAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	inputAccount := createRandomAccount(t)

	updateInfo := UpdateAccountParams{
		ID:      inputAccount.ID,
		Balance: util.RandomMoney(),
	}
	expectAccount, err := testQueries.UpdateAccount(context.Background(), updateInfo)

	require.NoError(t, err)
	require.NotEmpty(t, expectAccount)

	require.Equal(t, inputAccount.ID, expectAccount.ID)
	require.Equal(t, inputAccount.Owner, expectAccount.Owner)
	require.Equal(t, updateInfo.Balance, expectAccount.Balance)
	require.Equal(t, inputAccount.Currency, expectAccount.Currency)
	require.WithinDuration(t, inputAccount.CreatedAt, expectAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	inputAccount := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), inputAccount.ID)
	require.NoError(t, err)

	expect, err := testQueries.GetAccount(context.Background(), inputAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, expect)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
