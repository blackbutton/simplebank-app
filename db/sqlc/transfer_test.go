package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simplebank-app/util"
	"testing"
	"time"
)

func CreateRandomTransfer(t *testing.T) Transfer {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	return CreateRandomTransferWithAccount(account1, account2, t)
}

func CreateRandomTransferWithAccount(account1, account2 Account, t *testing.T) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, account1.ID)
	require.Equal(t, transfer.ToAccountID, account2.ID)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func DeleteRandomAccountTransfer(t *testing.T, fromAccount bool) {
	transfer1 := CreateRandomTransfer(t)
	var err error
	if fromAccount {
		err = testQueries.DeleteFromAccountTransfer(context.Background(), transfer1.FromAccountID)
	} else {
		err = testQueries.DeleteToAccountTransfer(context.Background(), transfer1.ToAccountID)
	}
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.Equal(t, err, sql.ErrNoRows)
	require.Empty(t, transfer2)
}

func TestDeleteFromAccountTransfer(t *testing.T) {
	DeleteRandomAccountTransfer(t, true)
}

func TestDeleteToAccountTransfer(t *testing.T) {
	DeleteRandomAccountTransfer(t, false)
}

func TestDeleteTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.Equal(t, err, sql.ErrNoRows)

	require.Empty(t, transfer2)
}

func ListAccountTransfer(t *testing.T, fromAccount bool) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	for i := 0; i < 10; i++ {
		CreateRandomTransferWithAccount(account1, account2, t)
	}
	var (
		transfers []Transfer
		err       error
	)
	if fromAccount {
		arg := ListFRomAccountTransfersParams{
			FromAccountID: account1.ID,
			Limit:         5,
			Offset:        5,
		}
		transfers, err = testQueries.ListFRomAccountTransfers(context.Background(), arg)
	} else {
		arg := ListToAccountTransfersParams{
			ToAccountID: account2.ID,
			Limit:       5,
			Offset:      5,
		}
		transfers, err = testQueries.ListToAccountTransfers(context.Background(), arg)
	}
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestListFRomAccountTransfers(t *testing.T) {
	ListAccountTransfer(t, true)
}

func TestToAccountTransfer(t *testing.T) {
	ListAccountTransfer(t, false)
}
