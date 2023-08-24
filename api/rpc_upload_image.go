package api

import (
	"bytes"
	"context"
	cloudflareGo "github.com/cloudflare/cloudflare-go"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
)

func (server *Server) UploadImage(ctx context.Context, request *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {
	violations := validateUploadImageRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

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

	uploadRequest := cloudflareGo.UploadImageParams{
		File: readCloser,
		Name: request.GetFilename(),
	}

	resourceContainer := cloudflareGo.ResourceContainer{
		Level:      cloudflareGo.AccountRouteLevel,
		Identifier: server.config.CloudflareAccountId,
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
	}

	return rsp, nil
}

func validateUploadImageRequest(req *pb.UploadImageRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateFilename(req.GetFilename()); err != nil {
		violations = append(violations, FieldViolation("filename", err))
	}

	if err := validator.ValidateImageTypeId(req.GetImageTypeId()); err != nil {
		violations = append(violations, FieldViolation("image_type_id", err))
	}

	return violations
}
