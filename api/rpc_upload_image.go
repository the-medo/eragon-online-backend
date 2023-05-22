package api

import (
	"bytes"
	"context"
	cloudflareGo "github.com/cloudflare/cloudflare-go"
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
