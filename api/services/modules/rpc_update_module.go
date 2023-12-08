package modules

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceModules) UpdateModule(ctx context.Context, req *pb.UpdateModuleRequest) (*pb.ViewModule, error) {
	violations := validateUpdateModuleRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	var needsEntityPermission []db.EntityType

	if req.DescriptionPostId != nil {
		needsEntityPermission = append(needsEntityPermission, db.EntityTypePost)
	}
	if req.AvatarImgId != nil || req.ThumbnailImgId != nil || req.HeaderImgId != nil {
		needsEntityPermission = append(needsEntityPermission, db.EntityTypeImage)
	}

	_, err := server.CheckModuleIdPermissions(ctx, req.GetModuleId(), &servicecore.ModulePermission{
		NeedsSuperAdmin:       true,
		NeedsEntityPermission: &needsEntityPermission,
	})

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to update module: %v", err)
	}

	arg := db.UpdateModuleParams{
		ID: req.GetModuleId(),
		AvatarImgID: sql.NullInt32{
			Int32: req.GetAvatarImgId(),
			Valid: req.AvatarImgId != nil,
		},
		ThumbnailImgID: sql.NullInt32{
			Int32: req.GetThumbnailImgId(),
			Valid: req.ThumbnailImgId != nil,
		},
		HeaderImgID: sql.NullInt32{
			Int32: req.GetHeaderImgId(),
			Valid: req.HeaderImgId != nil,
		},
		DescriptionPostID: sql.NullInt32{
			Int32: req.GetDescriptionPostId(),
			Valid: req.DescriptionPostId != nil,
		},
	}

	_, err = server.Store.UpdateModule(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update module: %v", err)
	}

	module, err := server.Store.GetModuleById(ctx, req.GetModuleId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve updated module: %v", err)
	}

	rsp := converters.ConvertViewModule(module)

	return rsp, nil
}

func validateUpdateModuleRequest(req *pb.UpdateModuleRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	if req.AvatarImgId != nil {
		if err := validator.ValidateImageId(req.GetAvatarImgId()); err != nil {
			violations = append(violations, e.FieldViolation("avatar_img_id", err))
		}
	}

	if req.ThumbnailImgId != nil {
		if err := validator.ValidateImageId(req.GetThumbnailImgId()); err != nil {
			violations = append(violations, e.FieldViolation("thumbnail_img_id", err))
		}
	}

	if req.HeaderImgId != nil {
		if err := validator.ValidateImageId(req.GetHeaderImgId()); err != nil {
			violations = append(violations, e.FieldViolation("header_img_id", err))
		}
	}

	if req.DescriptionPostId != nil {
		if err := validator.ValidatePostId(req.GetDescriptionPostId()); err != nil {
			violations = append(violations, e.FieldViolation("description_post_id", err))
		}
	}

	return violations
}
