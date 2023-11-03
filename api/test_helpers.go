package api

import (
	"bytes"
	"database/sql"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/util"
	"io/ioutil"
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

func RequireMatchUser(t *testing.T, user1 pb.User, user2 db.User) {
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	//require.Empty(t, user1.HashedPassword)
}

func RequireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
