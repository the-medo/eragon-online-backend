package testutils

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/util"
	"testing"
	"time"
)

func RandomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		ID:             util.RandomUserId(),
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
		ImgID: sql.NullInt32{
			Int32: util.RandomImgId(),
			Valid: true,
		},
		IsEmailVerified:   false,
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
	}
	return
}

func RandomViewUser(t *testing.T, user db.User) (viewUser db.ViewUser) {

	viewUser = db.ViewUser{
		ID:                user.ID,
		Username:          user.Username,
		HashedPassword:    user.HashedPassword,
		Email:             user.Email,
		ImgID:             user.ImgID,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		IsEmailVerified:   user.IsEmailVerified,
	}
	return
}
