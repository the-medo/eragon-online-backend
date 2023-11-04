package converters

import db "github.com/the-medo/talebound-backend/db/sqlc"

func ConvertViewUserToUser(user db.ViewUser) db.User {
	dbUser := db.User{
		ID:                 user.ID,
		Username:           user.Username,
		Email:              user.Email,
		ImgID:              user.ImgID,
		PasswordChangedAt:  user.PasswordChangedAt,
		CreatedAt:          user.CreatedAt,
		IsEmailVerified:    user.IsEmailVerified,
		IntroductionPostID: user.IntroductionPostID,
	}

	return dbUser
}
