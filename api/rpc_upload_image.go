package api

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	cloudflareGo "github.com/cloudflare/cloudflare-go"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
)

func (server *Server) UploadImage(ctx context.Context, request *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {
	cloudflare, err := cloudflareGo.NewWithAPIToken(server.config.CloudflareApiToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create cloudflare client: %v", err)
	}

	reader := bytes.NewReader(request.GetData())
	readCloser := io.NopCloser(reader)

	defer func(readCloser io.ReadCloser) {
		err := readCloser.Close()
		if err != nil {

		}
	}(readCloser)

	uploadRequest := cloudflareGo.ImageUploadRequest{
		File: readCloser,
		Name: request.GetFilename(),
	}

	img, err := cloudflare.UploadImage(ctx, server.config.CloudflareAccountId, uploadRequest)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload image: %v", err)
	}

	rsp := &pb.UploadImageResponse{
		Id:         img.ID,
		Filename:   img.Filename,
		Variants:   img.Variants,
		UploadedAt: timestamppb.New(img.Uploaded),
	}

	return rsp, nil
}

func (server *Server) UploadAndInsertToDb(ctx context.Context, data []byte, imgTypeId ImageTypeIds, filename string, userId int32) (*db.Image, error) {

	//upload image to cloudflare
	uploadRequest := &pb.UploadImageRequest{
		Filename: filename,
		Data:     data,
	}

	uploadImg, err := server.UploadImage(ctx, uploadRequest)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload user avatar: %v", err)
	}

	createImageParams, err := convertCloudflareImgToDb(server, ctx, uploadImg, imgTypeId, filename, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert Cloudflare image to DB: %v", err)
	}

	//insert img into DB "images" table
	dbImg, err := server.store.CreateImage(ctx, createImageParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert image into DB: %v", err)
	}

	return &dbImg, nil
}

func (server *Server) UploadUserAvatar(ctx context.Context, request *pb.UploadUserAvatarRequest) (*pb.UploadUserAvatarResponse, error) {
	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload.UserId != request.GetUserId() {
		return nil, status.Errorf(codes.PermissionDenied, "you are not allowed to update this user")
	}

	filename := fmt.Sprintf("avatar-%d", request.GetUserId())

	dbImg, err := server.UploadAndInsertToDb(ctx, request.GetData(), ImageTypeIdUserAvatar, filename, request.GetUserId())
	if err != nil {
		return nil, err
	}

	//update user imgID in DB
	_, err = server.store.UpdateUser(ctx, db.UpdateUserParams{
		ID: request.GetUserId(),
		ImgID: sql.NullInt32{
			Int32: dbImg.ID,
			Valid: true,
		},
	})

	if err != nil {
		return nil, err
	}

	return &pb.UploadUserAvatarResponse{
		UserId: request.GetUserId(),
		Image:  convertImage(dbImg),
	}, nil
}
