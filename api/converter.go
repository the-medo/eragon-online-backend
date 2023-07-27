package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/url"
	"path"
)

func convertImage(dbImage *db.Image) *pb.Image {
	pbImage := &pb.Image{
		Id:          dbImage.ID,
		ImgGuid:     dbImage.ImgGuid.UUID.String(),
		ImageTypeId: dbImage.ImageTypeID.Int32,
		Name:        dbImage.Name.String,
		Url:         dbImage.Url,
		BaseUrl:     dbImage.BaseUrl,
		CreatedAt:   timestamppb.New(dbImage.CreatedAt),
		UserId:      dbImage.UserID,
	}
	return pbImage
}

func convertUser(user db.User, img *pb.Image) *pb.User {
	pbUser := &pb.User{
		Id:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
		IsEmailVerified:   user.IsEmailVerified,
	}

	if user.ImgID.Valid == true {
		pbUser.ImgId = &user.ImgID.Int32
		pbUser.Img = img
	}

	return pbUser
}

func convertUserGetImage(server *Server, ctx context.Context, user db.User) *pb.User {
	pbUser := convertUser(user, nil)

	if user.ImgID.Valid == true {
		img, err := server.store.GetImageById(ctx, *pbUser.ImgId)
		if err != nil {
			return nil
		}
		pbUser.Img = convertImage(&img)
	}

	return pbUser
}

func convertViewUser(user db.ViewUser) *pb.ViewUser {
	pbUser := &pb.ViewUser{
		Id:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
		IsEmailVerified:   user.IsEmailVerified,
	}

	if user.ImgID.Valid == true {
		pbUser.ImgId = &user.ImgID.Int32
	}

	if user.AvatarImageUrl.Valid == true {
		pbUser.AvatarImageUrl = &user.AvatarImageUrl.String
	}

	if user.AvatarImageGuid.Valid == true {
		avatarImageGuid := user.AvatarImageGuid.UUID.String()
		pbUser.AvatarImageGuid = &avatarImageGuid
	}

	if user.IntroductionPostID.Valid == true {
		pbUser.IntroductionPostId = &user.IntroductionPostID.Int32
	}

	if user.IntroductionPostDeletedAt.Valid == true {
		pbUser.IntroductionPostDeletedAt = timestamppb.New(user.IntroductionPostDeletedAt.Time)
	}

	return pbUser
}

func convertViewUserToUser(user db.ViewUser) db.User {
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

func convertUserRowWithImage(user db.GetUsersRow) *pb.User {
	pbUser := &pb.User{
		Id:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
		IsEmailVerified:   user.IsEmailVerified,
	}

	if user.ImgID.Valid == true {
		pbUser.ImgId = &user.ImgID.Int32
		pbUser.Img = &pb.Image{
			Id:          user.ImgID.Int32,
			ImgGuid:     user.ImgGuid.UUID.String(),
			ImageTypeId: user.ImageTypeID.Int32,
			Name:        user.Name.String,
			Url:         user.Url.String,
			BaseUrl:     user.BaseUrl.String,
			CreatedAt:   timestamppb.New(user.CreatedAt_2.Time),
		}
	}

	return pbUser
}

func convertChatMessage(msg db.GetChatMessagesRow) *pb.ChatMessage {
	pbMessage := &pb.ChatMessage{
		Id:        msg.ChatID,
		UserId:    msg.UserID,
		Username:  msg.Username,
		Text:      msg.Text,
		CreatedAt: timestamppb.New(msg.CreatedAt),
	}

	return pbMessage
}

func convertEvaluationVote(evaluationVote db.EvaluationVote) *pb.EvaluationVote {
	pbEvaluationVote := &pb.EvaluationVote{
		EvaluationId: evaluationVote.EvaluationID,
		UserId:       evaluationVote.UserID,
		UserIdVoter:  evaluationVote.UserIDVoter,
		Value:        evaluationVote.Value,
		CreatedAt:    timestamppb.New(evaluationVote.CreatedAt),
	}

	return pbEvaluationVote
}

func convertCloudflareImgToDb(server *Server, ctx context.Context, uploadImg *pb.UploadImageResponse, imgTypeId ImageTypeIds, filename string, userId int32) (db.CreateImageParams, error) {
	nullImgTypeId := sql.NullInt32{
		Int32: int32(imgTypeId),
		Valid: true,
	}

	imageType, err := server.store.GetImageTypeById(ctx, int32(imgTypeId))
	if err != nil {
		return db.CreateImageParams{}, status.Errorf(codes.Internal, "failed to get image type: %v", err)
	}

	//format : https://imagedelivery.net/<account_id>/<image_id>/<variant_name>
	u, err := url.Parse(uploadImg.GetVariants()[0])
	if err != nil {
		return db.CreateImageParams{}, status.Errorf(codes.Internal, "failed to parse url: %v", err)
	}

	u.Path = path.Dir(u.Path) //removes last part of URL
	baseUrl := u.String()

	variantUrl := fmt.Sprintf("%s/%s", baseUrl, imageType.Variant)
	imgGuid := uuid.MustParse(uploadImg.GetId())

	rsp := db.CreateImageParams{
		ImgGuid: uuid.NullUUID{
			UUID:  imgGuid,
			Valid: true,
		},
		ImgTypeID: nullImgTypeId,
		Name: sql.NullString{
			String: fmt.Sprintf("%s_%s", filename, imgGuid.String()),
			Valid:  true,
		},
		Url:     variantUrl,
		BaseUrl: baseUrl,
		UserID:  userId,
	}

	return rsp, nil
}

func convertWorld(world db.ViewWorld) *pb.World {
	pbWorld := &pb.World{
		Id:             world.ID,
		Name:           world.Name,
		Public:         world.Public,
		CreatedAt:      timestamppb.New(world.CreatedAt),
		Description:    world.Description,
		ImageAvatar:    world.ImageAvatar.String,
		ImageThumbnail: world.ImageThumbnail.String,
		ImageHeader:    world.ImageHeader.String,
	}

	return pbWorld
}

func convertPostType(postType db.PostType) *pb.DataPostType {
	pbPostType := &pb.DataPostType{
		Id:         postType.ID,
		Name:       postType.Name,
		Draftable:  postType.Draftable,
		Privatable: postType.Privatable,
	}

	return pbPostType
}

func convertPostAndPostType(post db.Post, postType db.PostType) *pb.Post {
	pbPost := &pb.Post{
		Post: &pb.DataPost{
			Id:         post.ID,
			PostTypeId: post.PostTypeID,
			UserId:     post.UserID,
			Title:      post.Title,
			Content:    post.Content,
			CreatedAt:  timestamppb.New(post.CreatedAt),
			IsDraft:    post.IsDraft,
			IsPrivate:  post.IsPrivate,
		},
		PostType: &pb.DataPostType{
			Id:         postType.ID,
			Name:       postType.Name,
			Draftable:  postType.Draftable,
			Privatable: postType.Privatable,
		},
	}

	if post.DeletedAt.Valid == true {
		pbPost.Post.DeletedAt = timestamppb.New(post.DeletedAt.Time)
	}

	if post.LastUpdatedAt.Valid == true {
		pbPost.Post.LastUpdatedAt = timestamppb.New(post.LastUpdatedAt.Time)
	}

	if post.LastUpdatedUserID.Valid == true {
		pbPost.Post.LastUpdatedUserId = post.LastUpdatedUserID.Int32
	}

	return pbPost
}

func convertViewPost(viewPost db.ViewPost) *pb.Post {
	pbPost := &pb.Post{
		Post: &pb.DataPost{
			Id:         viewPost.ID,
			PostTypeId: viewPost.PostTypeID,
			UserId:     viewPost.UserID,
			Title:      viewPost.Title,
			Content:    viewPost.Content,
			CreatedAt:  timestamppb.New(viewPost.CreatedAt),
			IsDraft:    viewPost.IsDraft,
			IsPrivate:  viewPost.IsPrivate,
		},
		PostType: &pb.DataPostType{
			Id:         viewPost.PostTypeID,
			Name:       viewPost.PostTypeName,
			Draftable:  viewPost.PostTypeDraftable,
			Privatable: viewPost.PostTypePrivatable,
		},
	}

	if viewPost.DeletedAt.Valid == true {
		pbPost.Post.DeletedAt = timestamppb.New(viewPost.DeletedAt.Time)
	}

	if viewPost.LastUpdatedAt.Valid == true {
		pbPost.Post.LastUpdatedAt = timestamppb.New(viewPost.LastUpdatedAt.Time)
	}

	if viewPost.LastUpdatedUserID.Valid == true {
		pbPost.Post.LastUpdatedUserId = viewPost.LastUpdatedUserID.Int32
	}

	return pbPost
}

func convertHistoryPostWithoutContent(postHistory db.GetPostHistoryByPostIdRow, postType db.PostType) *pb.HistoryPost {
	pbHistoryPost := &pb.HistoryPost{
		Post: &pb.DataHistoryPost{
			Id:         postHistory.PostHistoryID,
			PostId:     postHistory.PostID,
			PostTypeId: postHistory.PostTypeID,
			UserId:     postHistory.UserID,
			Title:      postHistory.Title,
			CreatedAt:  timestamppb.New(postHistory.CreatedAt),
			IsDraft:    postHistory.IsDraft,
			IsPrivate:  postHistory.IsPrivate,
		},
		PostType: &pb.DataPostType{
			Id:         postType.ID,
			Name:       postType.Name,
			Draftable:  postType.Draftable,
			Privatable: postType.Privatable,
		},
	}

	if postHistory.DeletedAt.Valid == true {
		pbHistoryPost.Post.DeletedAt = timestamppb.New(postHistory.DeletedAt.Time)
	}

	if postHistory.LastUpdatedAt.Valid == true {
		pbHistoryPost.Post.LastUpdatedAt = timestamppb.New(postHistory.LastUpdatedAt.Time)
	}

	if postHistory.LastUpdatedUserID.Valid == true {
		pbHistoryPost.Post.LastUpdatedUserId = postHistory.LastUpdatedUserID.Int32
	}

	return pbHistoryPost
}

func convertHistoryPost(postHistory db.GetPostHistoryByIdRow, postType db.PostType) *pb.HistoryPost {
	pbHistoryPost := &pb.HistoryPost{
		Post: &pb.DataHistoryPost{
			Id:         postHistory.PostHistoryID,
			PostId:     postHistory.PostID,
			PostTypeId: postHistory.PostTypeID,
			UserId:     postHistory.UserID,
			Title:      postHistory.Title,
			Content:    postHistory.Content,
			CreatedAt:  timestamppb.New(postHistory.CreatedAt),
			IsDraft:    postHistory.IsDraft,
			IsPrivate:  postHistory.IsPrivate,
		},
		PostType: &pb.DataPostType{
			Id:         postType.ID,
			Name:       postType.Name,
			Draftable:  postType.Draftable,
			Privatable: postType.Privatable,
		},
	}

	if postHistory.DeletedAt.Valid == true {
		pbHistoryPost.Post.DeletedAt = timestamppb.New(postHistory.DeletedAt.Time)
	}

	if postHistory.LastUpdatedAt.Valid == true {
		pbHistoryPost.Post.LastUpdatedAt = timestamppb.New(postHistory.LastUpdatedAt.Time)
	}

	if postHistory.LastUpdatedUserID.Valid == true {
		pbHistoryPost.Post.LastUpdatedUserId = postHistory.LastUpdatedUserID.Int32
	}

	return pbHistoryPost
}
