package api

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateWorld(ctx context.Context, req *pb.CreateWorldRequest) (*pb.World, error) {
	violations := validateCreateWorldRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	arg := db.CreateWorldTxParams{
		CreateWorldParams: db.CreateWorldParams{
			Name:             req.GetName(),
			BasedOn:          req.GetBasedOn(),
			ShortDescription: req.GetShortDescription(),
		},
		UserId: authPayload.UserId,
	}

	txResult, err := server.store.CreateWorldTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "name already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create world: %s", err)
	}

	rsp := convertWorld(txResult)

	return rsp, nil
}

func (server *Server) UpdateWorld(ctx context.Context, req *pb.UpdateWorldRequest) (*pb.World, error) {
	violations := validateUpdateWorldRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	err := server.CheckWorldAdmin(ctx, req.GetWorldId(), true)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to update world: %v", err)
	}

	arg := db.UpdateWorldParams{
		WorldID: req.GetWorldId(),
		Name: sql.NullString{
			String: req.GetName(),
			Valid:  req.Name != nil,
		},
		BasedOn: sql.NullString{
			String: req.GetBasedOn(),
			Valid:  req.BasedOn != nil,
		},
		ShortDescription: sql.NullString{
			String: req.GetShortDescription(),
			Valid:  req.ShortDescription != nil,
		},
		Public: sql.NullBool{
			Bool:  req.GetPublic(),
			Valid: req.Public != nil,
		},
		DescriptionPostID: sql.NullInt32{
			Int32: req.GetDescriptionPostId(),
			Valid: req.DescriptionPostId != nil,
		},
	}

	changesMade := arg.Name.Valid || arg.BasedOn.Valid || arg.ShortDescription.Valid || arg.Public.Valid || arg.DescriptionPostID.Valid
	if changesMade {
		_, err = server.store.UpdateWorld(ctx, arg)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update world: %v", err)
		}
	}

	argImages := db.UpdateWorldImagesParams{
		WorldID: req.GetWorldId(),
		ThumbnailImgID: sql.NullInt32{
			Int32: req.GetImageThumbnailId(),
			Valid: req.ImageThumbnailId != nil,
		},
		HeaderImgID: sql.NullInt32{
			Int32: req.GetImageHeaderId(),
			Valid: req.ImageHeaderId != nil,
		},
		AvatarImgID: sql.NullInt32{
			Int32: req.GetImageAvatarId(),
			Valid: req.ImageAvatarId != nil,
		},
	}

	changesMade = argImages.ThumbnailImgID.Valid || argImages.HeaderImgID.Valid || argImages.AvatarImgID.Valid

	if changesMade {
		_, err = server.store.UpdateWorldImages(ctx, argImages)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update world: %v", err)
		}
	}

	world, err := server.store.GetWorldByID(ctx, req.GetWorldId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve updated world: %v", err)
	}

	rsp := convertWorld(world)

	return rsp, nil
}

func validateCreateWorldRequest(req *pb.CreateWorldRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateString(req.GetName(), 3, 64); err != nil {
		violations = append(violations, FieldViolation("name", err))
	}

	if err := validator.ValidateString(req.GetShortDescription(), 0, 1000); err != nil {
		violations = append(violations, FieldViolation("short_description", err))
	}

	if err := validator.ValidateString(req.GetBasedOn(), 0, 100); err != nil {
		violations = append(violations, FieldViolation("based_on", err))
	}

	return violations
}

func validateUpdateWorldRequest(req *pb.UpdateWorldRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if req.Name != nil {
		if err := validator.ValidateString(req.GetName(), 3, 64); err != nil {
			violations = append(violations, FieldViolation("name", err))
		}
	}

	if req.ShortDescription != nil {
		if err := validator.ValidateString(req.GetShortDescription(), 0, 1000); err != nil {
			violations = append(violations, FieldViolation("short_description", err))
		}
	}

	if req.BasedOn != nil {
		if err := validator.ValidateString(req.GetBasedOn(), 0, 100); err != nil {
			violations = append(violations, FieldViolation("based_on", err))
		}
	}

	return violations
}
