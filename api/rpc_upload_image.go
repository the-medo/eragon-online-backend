package api

import (
	"bytes"
	"context"
	cloudflareGo "github.com/cloudflare/cloudflare-go"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"image"
	_ "image/png" // Import this to decode png images
	"io"
)

func (server *Server) UploadImage(ctx context.Context, request *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {
	violations := validateUploadImageRequest(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	cloudflare, err := cloudflareGo.NewWithAPIToken(server.Config.CloudflareApiToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create cloudflare client: %v", err)
	}

	reader := bytes.NewReader(request.GetData())

	decodedImg, _, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to decode image: %v", err)
	}
	width := decodedImg.Width
	height := decodedImg.Height

	// Reset the reader to the beginning so it can be used again for uploading
	reader.Seek(0, 0)

	readCloser := io.NopCloser(reader)

	defer func(readCloser io.ReadCloser) {
		err := readCloser.Close()
		if err != nil {

		}
	}(readCloser)

	uploadRequest := cloudflareGo.UploadImageParams{
		File: readCloser,
		Name: request.GetFilename(),
	}

	resourceContainer := cloudflareGo.ResourceContainer{
		Level:      cloudflareGo.AccountRouteLevel,
		Identifier: server.Config.CloudflareAccountId,
		Type:       cloudflareGo.AccountType,
	}

	img, err := cloudflare.UploadImage(ctx, &resourceContainer, uploadRequest)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload image: %v", err)
	}

	rsp := &pb.UploadImageResponse{
		Id:         img.ID,
		Filename:   img.Filename,
		Variants:   img.Variants,
		UploadedAt: timestamppb.New(img.Uploaded),
		Width:      int32(width),  // assuming Width is int32 in your pb.UploadImageResponse
		Height:     int32(height), // assuming Height is int32 in your pb.UploadImageResponse
	}

	return rsp, nil
}

func validateUploadImageRequest(req *pb.UploadImageRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateFilename(req.GetFilename()); err != nil {
		violations = append(violations, e.FieldViolation("filename", err))
	}

	if err := validator.ValidateImageTypeId(req.GetImageTypeId()); err != nil {
		violations = append(violations, e.FieldViolation("image_type_id", err))
	}

	return violations
}
