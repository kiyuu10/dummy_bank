package db

import (
	"context"
	"testing"

	"github.com/kiyuu10/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	entryParams := CreateEntryParams{
		AccountID: account.ID,
		Amount:    10000,
	}

	entry, err := testQueries.CreateEntry(context.Background(), entryParams)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entryParams.AccountID, entry.AccountID)
	require.Equal(t, entryParams.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)
	require.NotEmpty(t, entry)

	entryGot, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryGot)

	require.Equal(t, entry.ID, entryGot.ID)
	require.Equal(t, entry.AccountID, entryGot.AccountID)
	require.Equal(t, entry.Amount, entryGot.Amount)
}

func TestEntryList(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
