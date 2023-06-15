package db

import (
	"context"
	"simple_bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func getAccountsFortTransfer(t *testing.T) (Account, Account) {
	arg := ListAccountParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, accounts)
	var account1 Account
	var account2 Account
	if len(accounts) > 2 {
		account1 = accounts[0]
		account2 = accounts[2]
	} else {
		var accErr1 error
		var accErr2 error
		account1, accErr1 = CreateRandomAccount(t)
		account2, accErr2 = CreateRandomAccount(t)
		require.NoError(t, accErr1)
		require.NoError(t, accErr2)
	}

	return account1, account2
}

func CreateRandomTransfer(t *testing.T) (Transfer, error) {

	account1, account2 := getAccountsFortTransfer(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	return transfer, err
}

func CreateRandomTransferForAccount(t *testing.T, from Account, to Account) {

	arg := CreateTransferParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Amount:        util.RandomMoney(),
	}

	for i := 0; i < 10; i++ {
		testQueries.CreateTransfer(context.Background(), arg)
	}

}

func TestCreateTransfer(t *testing.T) {

	transfer, err := CreateRandomTransfer(t)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NotZero(t, transfer.FromAccountID)
	require.NotZero(t, transfer.ToAccountID)
	require.NotZero(t, transfer.Amount)

}

func TestGetTransfer(t *testing.T) {
	transfer, err1 := CreateRandomTransfer(t)
	require.NoError(t, err1)

	transfer2, err2 := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err2)
	require.Equal(t, transfer.ID, transfer2.ID)
}

func TestListTransfers(t *testing.T) {
	account1, account2 := getAccountsFortTransfer(t)
	CreateRandomTransferForAccount(t, account1, account2)

	arg2 := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg2)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}

}
