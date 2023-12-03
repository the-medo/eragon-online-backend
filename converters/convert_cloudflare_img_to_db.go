package converters

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/constants"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
	"path"
)

func ConvertCloudflareImgToDb(server *servicecore.ServiceCore, ctx context.Context, uploadImg *pb.UploadImageResponse, imgTypeId constants.ImageTypeIds, filename string, userId int32) (db.CreateImageParams, error) {
	nullImgTypeId := sql.NullInt32{
		Int32: int32(imgTypeId),
		Valid: true,
	}

	imageType, err := server.Store.GetImageTypeById(ctx, int32(imgTypeId))
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
		Width:   uploadImg.GetWidth(),
		Height:  uploadImg.GetHeight(),
	}

	return rsp, nil
}
