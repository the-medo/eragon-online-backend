package images

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/the-medo/talebound-backend/api"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

/*
UploadDefaultImage
Uploads an image to Cloudflare and inserts it into the DB with image type Default
  - 1. filename: {request.filename}-{userId}
  - 2. upload to cloudflare => get cloudflareId
  - 3. insert into DB `{request.filename}-{userId}_{cloudflareId}`
*/
func (server *api.Server) UploadDefaultImage(ctx context.Context, request *pb.UploadImageRequest) (*pb.Image, error) {
	violations := validateUploadDefaultImage(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	filename := fmt.Sprintf("%s-%d", request.GetFilename(), authPayload.UserId)

	log.Info().Int32("ImageTypeId", request.GetImageTypeId()).Str("filename", request.GetFilename()).Msgf("Uploading default image: %s", filename)

	dbImg, err := server.UploadAndInsertToDb(ctx, request.GetData(), api.ImageTypeIds(request.GetImageTypeId()), filename, authPayload.UserId)
	if err != nil {
		return nil, err
	}

	return api.ConvertImage(dbImg), nil
}

func validateUploadDefaultImage(req *pb.UploadImageRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateFilename(req.GetFilename()); err != nil {
		violations = append(violations, e.FieldViolation("filename", err))
	}

	if err := validator.ValidateImageTypeId(req.GetImageTypeId()); err != nil {
		violations = append(violations, e.FieldViolation("image_type_id", err))
	}

	return violations
}
