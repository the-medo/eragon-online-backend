package db

import (
	"context"
	"github.com/the-medo/talebound-backend/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ResetPasswordRequestTxParams contains the input parameters of the reset transaction
type ResetPasswordRequestTxParams struct {
	Email       string
	AfterCreate func(user User) error
}

// ResetPasswordRequestTxResult contains the result of the reset transaction
type ResetPasswordRequestTxResult struct {
	User User `json:"user"`
}

// ResetPasswordRequestTx performs a reset password verification.
func (store *SQLStore) ResetPasswordRequestTx(ctx context.Context, arg ResetPasswordRequestTxParams) (ResetPasswordRequestTxResult, error) {
	var result ResetPasswordRequestTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.GetUserByEmail(ctx, arg.Email)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to get user: %s", err)
		}

		_, err = q.AddUserPasswordReset(ctx, AddUserPasswordResetParams{
			UserID: result.User.ID,
			Code:   util.RandomString(64),
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to add 'user password reset' row into DB: %s", err)
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
