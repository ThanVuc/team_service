package usecase

import (
	"context"
	appdto "team_service/internal/application/common/dto"
	apphelper "team_service/internal/application/common/helper"
	istore "team_service/internal/application/common/interface/store"
	appvalidation "team_service/internal/application/common/validation"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/proto/common"

	"github.com/thanvuc/go-core-lib/log"
	"github.com/wagslane/go-rabbitmq"
)

type (
	GroupUseCase interface {
		CreateGroup(ctx context.Context, req *appdto.CreateGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError)
		Ping(ctx context.Context, req *common.EmptyRequest) (*common.EmptyResponse, errorbase.AppError)
		GetGroupRequest(ctx context.Context, req *appdto.GetGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError)
		ListGroups(ctx context.Context, req *appdto.ListGroupsRequest) (*appdto.BaseResponse[appdto.ListGroupsResponse], errorbase.AppError)
		GetSimpleUserByGroupID(ctx context.Context, req *appdto.ListMembersRequest) (*appdto.BaseResponse[[]appdto.SimpleUserResponse], errorbase.AppError)
		UpdateGroup(ctx context.Context, req *appdto.UpdateGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError)
		DeleteGroup(ctx context.Context, req *appdto.DeleteGroupRequest) (*appdto.BaseResponse[appdto.DeleteGroupResponse], errorbase.AppError)
		GetListGroupMembers(ctx context.Context, req *appdto.ListMembersRequest) (*appdto.BaseResponse[appdto.ListMembersResponse], errorbase.AppError)
		UpdateMemberRole(ctx context.Context, req *appdto.UpdateMemberRoleRequest) (*appdto.BaseResponse[appdto.MemberResponse], errorbase.AppError)
		RemoveMember(ctx context.Context, req *appdto.RemoveMemberRequest) (*appdto.BaseResponse[appdto.RemoveMemberResponse], errorbase.AppError)
		CreateInvite(ctx context.Context, req *appdto.CreateInviteRequest) (*appdto.BaseResponse[appdto.InviteResponse], errorbase.AppError)
	}

	SprintUseCase interface {
		CreateSprint(ctx context.Context, req *appdto.CreateSprintRequest) (*appdto.BaseResponse[appdto.SprintResponse], errorbase.AppError)
		GetSprint(ctx context.Context, req *appdto.GetSprintRequest) (*appdto.BaseResponse[appdto.SprintResponse], errorbase.AppError)
		GetSimpleSprints(ctx context.Context, req *appdto.ListSprintsRequest) (*appdto.BaseResponse[[]appdto.SimpleSprintResponse], errorbase.AppError)
		ListSprints(ctx context.Context, req *appdto.ListSprintsRequest) (*appdto.BaseResponse[appdto.ListSprintsResponse], errorbase.AppError)
		UpdateSprint(ctx context.Context, req *appdto.UpdateSprintRequest) (*appdto.BaseResponse[appdto.SprintResponse], errorbase.AppError)
		UpdateSprintStatus(ctx context.Context, req *appdto.UpdateSprintStatusRequest) (*appdto.BaseResponse[appdto.UpdateSprintStatusResponse], errorbase.AppError)
		DeleteSprint(ctx context.Context, req *appdto.DeleteSprintRequest) (*appdto.BaseResponse[appdto.DeleteSprintResponse], errorbase.AppError)
		ExportSprint(ctx context.Context, req *appdto.ExportSprintRequest) (*appdto.BaseResponse[appdto.ExportSprintResponse], errorbase.AppError)
	}

	WorkUseCase interface {
		CreateWork(ctx context.Context, req *appdto.CreateWorkRequest) (*appdto.BaseResponse[appdto.WorkResponse], errorbase.AppError)
		GetWork(ctx context.Context, req *appdto.GetWorkRequest) (*appdto.BaseResponse[appdto.WorkResponse], errorbase.AppError)
		ListWorks(ctx context.Context, req *appdto.ListWorksRequest) (*appdto.BaseResponse[appdto.ListWorksResponse], errorbase.AppError)
		UpdateWork(ctx context.Context, req *appdto.UpdateWorkRequest) (*appdto.BaseResponse[appdto.WorkResponse], errorbase.AppError)
		DeleteWork(ctx context.Context, req *appdto.DeleteWorkRequest) (*appdto.BaseResponse[appdto.DeleteWorkResponse], errorbase.AppError)

		CreateChecklistItem(ctx context.Context, req *appdto.CreateChecklistItemRequest) (*appdto.BaseResponse[appdto.ChecklistItemResponse], errorbase.AppError)
		UpdateChecklistItem(ctx context.Context, req *appdto.UpdateChecklistItemRequest) (*appdto.BaseResponse[appdto.ChecklistItemResponse], errorbase.AppError)
		DeleteChecklistItem(ctx context.Context, req *appdto.DeleteChecklistItemRequest) (*appdto.BaseResponse[appdto.ChecklistItemResponse], errorbase.AppError)

		CreateComment(ctx context.Context, req *appdto.CreateCommentRequest) (*appdto.BaseResponse[appdto.CommentListResponse], errorbase.AppError)
		UpdateComment(ctx context.Context, req *appdto.UpdateCommentRequest) (*appdto.BaseResponse[appdto.CommentListResponse], errorbase.AppError)
		DeleteComment(ctx context.Context, req *appdto.DeleteCommentRequest) (*appdto.BaseResponse[appdto.CommentListResponse], errorbase.AppError)
	}

	UserUseCase interface {
		SyncUserData(ctx context.Context) func(d rabbitmq.Delivery) rabbitmq.Action
	}
)

func NewGroupUseCase(
	store istore.Store,
	validator *appvalidation.GroupValidator,
	authHelper *apphelper.AuthHelper,
	notificationHelper *apphelper.NotificationHelper,
) GroupUseCase {
	return &groupUseCase{
		store:              store,
		groupRepo:          store.GroupRepository(),
		userRepo:           store.UserRepository(),
		validator:          validator,
		notificationHelper: notificationHelper,
	}
}

func NewSprintUseCase(
	store istore.Store,
	validator *appvalidation.SprintValidator,
	authHelper *apphelper.AuthHelper,
) SprintUseCase {
	return &sprintUseCase{
		store:              store,
		sprintRepo:         store.SprintRepository(),
		workRepo:           store.WorkRepository(),
		userRepo:           store.UserRepository(),
		validator:          validator,
		authHelper:         authHelper,
		sprintExportHelper: apphelper.NewSprintExportHelper(),
	}
}

func NewWorkUseCase(
	store istore.Store,
	validator *appvalidation.WorkValidator,
	authHelper *apphelper.AuthHelper,
) WorkUseCase {
	return &workUseCase{
		store:      store,
		workRepo:   store.WorkRepository(),
		validator:  validator,
		authHelper: authHelper,
	}
}

func NewUserUseCase(
	store istore.Store,
	logger log.LoggerV2,
) UserUseCase {
	return &userUseCase{
		store:  store,
		logger: logger,
	}
}
