package db

import (
	"context"
	"database/sql"
	"simple_bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) (Account, error) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	return account, err
}

func TestCreateAccount(t *testing.T) {
	account, err := CreateRandomAccount(t)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.Owner)
	require.NotZero(t, account.Balance)
	require.NotZero(t, account.Currency)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	account1, err1 := CreateRandomAccount(t)
	require.NoError(t, err1)

	account2, err2 := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err2)
	require.Equal(t, account1.ID, account2.ID)
}

func TestUpdateAccount(t *testing.T) {
	account, err := CreateRandomAccount(t)
	require.NoError(t, err)

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}

	account1, err1 := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err1)

	require.Equal(t, account.ID, account1.ID)
	require.Equal(t, arg.Balance, account1.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account, err := CreateRandomAccount(t)
	require.NoError(t, err)

	testQueries.DeleteAccount(context.Background(), account.ID)
	account1, err1 := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err1)
	require.EqualError(t, err1, sql.ErrNoRows.Error())
	require.Empty(t, account1)

}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
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
