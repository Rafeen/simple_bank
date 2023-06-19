package db

import (
	"context"
	"database/sql"
	"fmt"
)

type (

	// Store provied all functions to execure db dueries and transactions
	Store struct {
		*Queries
		db *sql.DB
	}

	// TransferTxParams contains the input  parameters of the transfer transaction
	TransferTxParams struct {
		FromAccountID int64 `json:"from_account_id"`
		ToAccountID   int64 `json:"to_account_id"`
		Amount        int64 `json:"amount"`
	}

	// TransferTxResult contains the input  parameters of the transfer transaction
	TransferTxResult struct {
		Transfer    Transfer `json:"transfer"`
		FromAccount Account  `json:"from_account"`
		ToAccount   Account  `json:"to_account"`
		Amount      Account  `json:"amount"`
		FromEntry   Entry    `json:"from_entry"`
		ToEntry     Entry    `json:"to_entry"`
	}
)

var txKey = struct{}{}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rbErr: %v", err, rbErr)
		}
	}

	return tx.Commit()
}

// TransferTx performs a money transfer from one account to the other
// it creates a transfer record, add account entries and update account's balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	// Start Transaction
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)
		// first query
		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		// second query
		fmt.Println(txName, "create from entry")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "create to entry")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// third query
		fmt.Println(txName, "get from account")
		fromAccount, err := q.GetAccountForUpdate(context.Background(), arg.FromAccountID)

		if err != nil {
			return err
		}
		fmt.Println(txName, "update from account")
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     fromAccount.ID,
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "get to account")
		to_account, err := q.GetAccountForUpdate(context.Background(), arg.ToAccountID)

		if err != nil {
			return err
		}

		fmt.Println(txName, "update to account")
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     to_account.ID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
