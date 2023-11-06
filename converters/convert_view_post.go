package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//TODO: maybe, once generics will help... -_-
// complete copy of ConvertGetPostsByModuleRow, but with different name

func ConvertViewPost(viewPost db.ViewPost) *pb.ViewPost {
	pbPost := &pb.ViewPost{
		Id:        viewPost.ID,
		UserId:    viewPost.UserID,
		Title:     viewPost.Title,
		Content:   viewPost.Content,
		CreatedAt: timestamppb.New(viewPost.CreatedAt),
		IsDraft:   viewPost.IsDraft,
		IsPrivate: viewPost.IsPrivate,
		Tags:      viewPost.Tags,
	}

	if viewPost.DeletedAt.Valid == true {
		pbPost.DeletedAt = timestamppb.New(viewPost.DeletedAt.Time)
	}

	if viewPost.LastUpdatedAt.Valid == true {
		pbPost.LastUpdatedAt = timestamppb.New(viewPost.LastUpdatedAt.Time)
	}

	if viewPost.LastUpdatedUserID.Valid == true {
		pbPost.LastUpdatedUserId = viewPost.LastUpdatedUserID.Int32
	}

	if viewPost.Description.Valid == true {
		pbPost.Description = viewPost.Description.String
	}

	if viewPost.ThumbnailImgID.Valid == true {
		pbPost.ImageThumbnailId = viewPost.ThumbnailImgID.Int32
	}

	if viewPost.ThumbnailImgUrl.Valid == true {
		pbPost.ImageThumbnailUrl = viewPost.ThumbnailImgUrl.String
	}

	if viewPost.EntityID.Valid == true {
		pbPost.EntityId = &viewPost.EntityID.Int32
	}

	if viewPost.ModuleID.Valid == true {
		pbPost.ModuleId = &viewPost.ModuleID.Int32
	}

	if viewPost.ModuleTypeID.Valid == true {
		pbPost.ModuleId = &viewPost.ModuleTypeID.Int32
	}

	if viewPost.ModuleType.Valid == true {
		convertedModuleType := ConvertModuleTypeToPB(viewPost.ModuleType.ModuleType)
		pbPost.ModuleType = &convertedModuleType
	}

	return pbPost
}
