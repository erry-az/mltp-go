package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/erry-az/mltp-go/db/query"
	"github.com/labstack/echo/v4"
)

func (s *Service) Transfer(ctx context.Context, senderUsername, targetUsername string, amount uint64) error {
	sender, err := s.queries.GetUserByUsername(ctx, senderUsername)
	if err != nil {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, "sender not found")
	}

	if uint64(sender.Balance) < amount {
		return echo.NewHTTPError(http.StatusBadRequest, "insufficient balance")
	}

	target, err := s.queries.GetUserByUsername(ctx, targetUsername)
	if err != nil {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, "destination user not found")
	}

	tx, err := s.pool.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	q := s.queries.WithTx(tx)

	senderUpdate, err := q.UpdateBalance(ctx, query.UpdateBalanceParams{
		ID:      sender.ID,
		Balance: int64(amount) * -1,
	})
	if err != nil {
		return err
	}

	if sender.Balance == senderUpdate.Balance {
		return echo.NewHTTPError(http.StatusBadRequest, "insufficient balance")
	}

	_, err = q.CreateTransaction(ctx, query.CreateTransactionParams{
		UserID: sender.ID,
		Amount: int64(amount),
		Name:   query.TransactionNameTransfer,
		Type:   query.TransactionTypeDebit,
	})
	if err != nil {
		return err
	}

	targetUpdate, err := q.UpdateBalance(ctx, query.UpdateBalanceParams{
		ID:      target.ID,
		Balance: int64(amount),
	})
	if err != nil {
		return err
	}

	if target.Balance == targetUpdate.Balance {
		return echo.NewHTTPError(http.StatusBadRequest, "insufficient balance")
	}

	_, err = q.CreateTransaction(ctx, query.CreateTransactionParams{
		UserID: target.ID,
		Amount: int64(amount),
		Name:   query.TransactionNameTransfer,
		Type:   query.TransactionTypeCredit,
	})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) TopUp(ctx context.Context, username string, amount uint64) error {
	user, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	tx, err := s.pool.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	q := s.queries.WithTx(tx)

	_, err = q.UpdateBalance(ctx, query.UpdateBalanceParams{
		ID:      user.ID,
		Balance: int64(amount),
	})
	if err != nil {
		return err
	}

	_, err = q.CreateTransaction(ctx, query.CreateTransactionParams{
		UserID: user.ID,
		Amount: int64(amount),
		Name:   query.TransactionNameTopUp,
		Type:   query.TransactionTypeCredit,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) TopTransactionByUsername(ctx context.Context, username string) ([]query.TopTransactionByUsernameRow, error) {
	var limit = 10

	return s.queries.TopTransactionByUsername(ctx, query.TopTransactionByUsernameParams{
		Username: username,
		Limit:    int32(limit),
	})
}

func (s *Service) TopSummaryTransactions() ([]query.TopSummaryTransactionsRow, error) {
	var limit = 10

	return s.queries.TopSummaryTransactions(context.Background(), int32(limit))
}
