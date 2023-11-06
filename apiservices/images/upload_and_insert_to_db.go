package images

import (
	"context"
	"github.com/the-medo/talebound-backend/api/constants"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UploadAndInsertToDb(service *servicecore.ServiceCore, ctx context.Context, data []byte, imgTypeId constants.ImageTypeIds, filename string, userId int32) (*db.Image, error) {

	//upload image to cloudflare
	uploadRequest := &pb.UploadImageRequest{
		Filename:    filename,
		Data:        data,
		ImageTypeId: int32(imgTypeId),
	}

	uploadImg, err := UploadImage(service, ctx, uploadRequest)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload image: %v", err)
	}

	createImageParams, err := converters.ConvertCloudflareImgToDb(service, ctx, uploadImg, imgTypeId, filename, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert Cloudflare image to DB: %v", err)
	}

	//insert img into DB "images" table
	dbImg, err := service.Store.CreateImage(ctx, createImageParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert image into DB: %v", err)
	}

	return &dbImg, nil
}
