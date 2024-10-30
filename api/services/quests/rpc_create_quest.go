package quests

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceQuests) CreateQuest(ctx context.Context, req *pb.CreateQuestRequest) (*pb.CreateQuestResponse, error) {
	violations := validateCreateQuestRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	arg := db.CreateQuestTxParams{
		CreateQuestParams: db.CreateQuestParams{
			Name:             req.GetName(),
			ShortDescription: req.GetShortDescription(),
			WorldID: sql.NullInt32{
				Int32: req.WorldId,
				Valid: true,
			},
			SystemID: sql.NullInt32{
				Int32: req.SystemId,
				Valid: true,
			},
		},
		UserId: authPayload.UserId,
	}

	txResult, err := server.Store.CreateQuestTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "name already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create quest: %s", err)
	}

	rsp := &pb.CreateQuestResponse{
		Quest:  converters.ConvertQuest(*txResult.Quest),
		Module: converters.ConvertViewModule(*txResult.Module),
	}

	return rsp, nil
}

func validateCreateQuestRequest(req *pb.CreateQuestRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetSystemId()); err != nil {
		violations = append(violations, e.FieldViolation("system_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	if err := validator.ValidateModuleName(req.GetName()); err != nil {
		violations = append(violations, e.FieldViolation("name", err))
	}

	if err := validator.ValidateModuleShortDescription(req.GetShortDescription()); err != nil {
		violations = append(violations, e.FieldViolation("short_description", err))
	}

	return violations
}
