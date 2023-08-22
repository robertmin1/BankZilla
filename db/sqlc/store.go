package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TranferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store{
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

// executes a function withing a DB transaction
func(store *SQLStore) execTC(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v Rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount 		  int64 `json:"amount"`
}

type TranferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"form_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry	`json:"form_entry"`
	ToEntry Entry	`json:"to_entry"`
}


func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TranferTxResult, error) {
	var result TranferTxResult
	fromacc := sql.NullInt64 {
		Int64: arg.FromAccountID,
		Valid: true,
	}
	toacc := sql.NullInt64 {
		Int64: arg.ToAccountID,
		Valid: true,
	}
	
	err := store.execTC(ctx, func (q *Queries) error {
		var err error

		result.Transfer, err =  q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: fromacc,
			ToAccountID: toacc,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: fromacc,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}


		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: toacc,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		//TODO: update account balances
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount,arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		}else{
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount,arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}
		

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	AccountID1 int64,
	amount1 int64,
	AccountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID: AccountID1,
	})
	if err != nil {
		return 
	}
	
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID: AccountID2,
	})
	if err != nil {
		return 
	}

	return
}