package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	arg := TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        10,
	}

	n := 10
	results := make(chan TransferTxResult)
	errs := make(chan error)
	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), arg)
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, result.Transfer.FromAccountID, arg.FromAccountID)
		require.Equal(t, result.Transfer.ToAccountID, arg.ToAccountID)
		require.Equal(t, result.Transfer.Amount, arg.Amount)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, result.FromEntry.AccountID, arg.FromAccountID)
		require.Equal(t, result.FromEntry.Amount, -arg.Amount)
		require.NotZero(t, result.FromEntry.CreatedAt)
		_, err = testStore.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, result.ToEntry.AccountID, arg.ToAccountID)
		require.Equal(t, result.ToEntry.Amount, arg.Amount)
		require.NotZero(t, result.ToEntry.CreatedAt)
		_, err = testStore.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, arg.FromAccountID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, arg.ToAccountID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%arg.Amount == 0)

		k := int(diff1 / arg.Amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updateAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*arg.Amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*arg.Amount, updateAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	n := 10
	amount := 10
	errs := make(chan error)
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        int64(amount),
			})
			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updateAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
