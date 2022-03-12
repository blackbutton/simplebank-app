package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simplebank-app/util"
	"testing"
	"time"
)

func CreateRandomEntry(t *testing.T) Entry {
	account := CreateRandomAccount(t)
	return CreateRandomEntryWithAccount(account, t)
}

func CreateRandomEntryWithAccount(account Account, t *testing.T) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	require.Equal(t, entry.Amount, arg.Amount)
	require.Equal(t, entry.AccountID, arg.AccountID)
	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.Equal(t, err, sql.ErrNoRows)
	require.Empty(t, entry2)
}

func TestDeleteAccountEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	err := testQueries.DeleteAccountEntry(context.Background(), entry1.AccountID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.Equal(t, err, sql.ErrNoRows)
	require.Empty(t, entry2)
}

func TestListEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	for i := 0; i < 10; i++ {
		CreateRandomEntryWithAccount(account, t)
	}
	arg := ListEntryParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
