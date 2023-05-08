package db

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// ResetPasswordVerifyTxParams contains the input parameters of the reset transaction
type ResetPasswordVerifyTxParams struct {
	UserId      int32
	Code        string
	NewPassword string
}

// ResetPasswordVerifyTxResult contains the result of the reset transaction
type ResetPasswordVerifyTxResult struct {
	User User
}

// ResetPasswordVerifyTx performs a reset password verification.
func (store *SQLStore) ResetPasswordVerifyTx(ctx context.Context, arg ResetPasswordVerifyTxParams) (ResetPasswordVerifyTxResult, error) {
	var result ResetPasswordVerifyTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.GetUserPasswordReset(ctx, GetUserPasswordResetParams{
			UserID: arg.UserId,
			Code:   arg.Code,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to get user password reset: %s", err)
		}

		hashedPassword, err := util.HashPassword(arg.NewPassword)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			ID: arg.UserId,
			HashedPassword: sql.NullString{
				String: hashedPassword,
				Valid:  true,
			},
			PasswordChangedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		})

		err = q.DeleteUserPasswordReset(ctx, DeleteUserPasswordResetParams{
			UserID: arg.UserId,
			Code:   arg.Code,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to delete 'user password reset': %s", err)
		}

		return err
	})

	return result, err
}
