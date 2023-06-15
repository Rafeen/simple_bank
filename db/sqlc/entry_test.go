package db

import (
	"context"
	"simple_bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func getAccountsFortEntry(t *testing.T) Account {
	arg := ListAccountParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	var account Account
	if len(accounts) > 0 {
		account = accounts[0]
	} else {
		var accErr error
		account, accErr = CreateRandomAccount(t)
		require.NoError(t, accErr)
	}

	return account
}

func CreateRandomEntry(t *testing.T) (Entry, error) {
	account := getAccountsFortEntry(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	return entry, err
}

func CreateRandomEntryForAccount(t *testing.T) Account {
	account := getAccountsFortEntry(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	for i := 0; i < 10; i++ {
		testQueries.CreateEntry(context.Background(), arg)
	}

	return account

}

func TestCreateEntry(t *testing.T) {
	entry, err := CreateRandomEntry(t)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.NotZero(t, entry.AccountID)
	require.NotZero(t, entry.Amount)

}

func TestGetEntry(t *testing.T) {
	entry1, err1 := CreateRandomEntry(t)
	require.NoError(t, err1)

	entry2, err2 := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err2)
	require.Equal(t, entry1.ID, entry2.ID)
}

func TestListEntry(t *testing.T) {

	account := CreateRandomEntryForAccount(t)
	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
