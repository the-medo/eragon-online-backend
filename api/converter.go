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
	baseUrl := path.Dir(uploadImg.GetVariants()[0]) //removes last part of URL
	url := path.Join(baseUrl, "/", string(imageType.Variant))
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
		Url:     url,
		BaseUrl: baseUrl,
		UserID:  userId,
	}

	return rsp, nil
}
