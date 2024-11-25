package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		pool:    pool,
		Queries: New(pool),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Begin a transaction using the pool
	tx, err := store.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}


func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		var err error
		var createTransferParams CreateTransferParams
		createTransferParams.FromAccountID = arg.FromAccountID
		createTransferParams.ToAccountID = arg.ToAccountID
		createTransferParams.Amount = arg.Amount

		result.Transfer, err = q.CreateTransfer(ctx, createTransferParams)
		if err != nil {
			return err
		}


		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//** This is the way to make sure the transaction doesnt have a deadlock
		//* this make sure the transaction run in specific order
		//** update begin >> account 1 than account 2 >> commit than next(concurrent) transaction should happen in the same order

		//*! this will create deadlock
		//*! tx1 begin >> update acount 1 than account 2 
		//*! tx2 begin >> update account 2 thean account 1

		//*! Since if we see concurrently we will observe that account 2 has lock when tx1 is happning and the same time if tx2 start and tries to update account 2 since account 2 is locked and not able to get free to which results in deadlock 

		if arg.FromAccountID < arg.ToAccountID{

			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
			
			// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID: arg.FromAccountID,
			// 	Amount: -arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }


			// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:      arg.ToAccountID,
			// 	Amount: arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
		}else{

			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
			
			// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:      arg.ToAccountID,
			// 	Amount: arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
			
			// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID: arg.FromAccountID,
			// 	Amount: -arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }

		}

		return nil
	})

	return result, err

}

func addMoney(
	ctx context.Context, 
	q *Queries, 
	accountID1 int64,
	amount1 int64, 
	accountID2 int64, 
	amount2 int64)(account1 Account, account2 Account, err error){

		account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: accountID1,
			Amount: amount1,
		})
		if err != nil{
			return
		}

		account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: accountID2,
			Amount: amount2,
		})
		
		return

	}
