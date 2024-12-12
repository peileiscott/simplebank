package db

import "context"

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToAccount   Account  `json:"to_account"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
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

		updateAccountBalance := func(accountID, amount int64) (Account, error) {
			return q.IncreaseAccountBalance(ctx, IncreaseAccountBalanceParams{
				ID:     accountID,
				Amount: amount,
			})
		}

		// Avoid potential deadlocks. See TestSQLStore_TransferTxDeadlock in test_store.go
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = updateAccountBalance(arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}

			result.ToAccount, err = updateAccountBalance(arg.ToAccountID, arg.Amount)
			return err
		} else {
			result.ToAccount, err = updateAccountBalance(arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}

			result.FromAccount, err = updateAccountBalance(arg.FromAccountID, -arg.Amount)
			return err
		}
	})

	return result, err
}
