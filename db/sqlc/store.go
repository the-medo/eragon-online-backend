package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	CreateWorldTx(ctx context.Context, arg CreateWorldTxParams) (CreateWorldTxResult, error)
	CreateSystemTx(ctx context.Context, arg CreateSystemTxParams) (CreateSystemTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
	ResetPasswordRequestTx(ctx context.Context, arg ResetPasswordRequestTxParams) (ResetPasswordRequestTxResult, error)
	ResetPasswordVerifyTx(ctx context.Context, arg ResetPasswordVerifyTxParams) (ResetPasswordVerifyTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
