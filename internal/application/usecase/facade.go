package usecase

import (
	"context"
	appdto "team_service/internal/application/common/dto"
	apphelper "team_service/internal/application/common/helper"
	icacherepository "team_service/internal/application/common/interface/cacherepository"
	istore "team_service/internal/application/common/interface/store"
	appvalidation "team_service/internal/application/common/validation"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/proto/common"

	"github.com/thanvuc/go-core-lib/log"
	"github.com/thanvuc/go-core-lib/storage"
	"github.com/wagslane/go-rabbitmq"
)

type (
	GroupUseCase interface {
		CreateGroup(ctx context.Context, req *appdto.CreateGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError)
		Ping(ctx context.Context, req *common.EmptyRequest) (*common.EmptyResponse, errorbase.AppError)
		ListGroups(ctx context.Context, req *appdto.ListGroupsRequest) (*appdto.BaseResponse[appdto.ListGroupsResponse], errorbase.AppError)
		GetGroup(ctx context.Context, req *appdto.GetGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError)
		GetSimpleUserByGroupID(ctx context.Context, req *appdto.ListMembersRequest) (*appdto.BaseResponse[[]appdto.SimpleUserResponse], errorbase.AppError)
		UpdateGroup(ctx context.Context, req *appdto.UpdateGroupRequest) (*appdto.BaseResponse[appdto.GroupResponse], errorbase.AppError)
		DeleteGroup(ctx context.Context, req *appdto.DeleteGroupRequest) (*appdto.BaseResponse[appdto.DeleteGroupResponse], errorbase.AppError)
		GetListGroupMembers(ctx context.Context, req *appdto.ListMembersRequest) (*appdto.BaseResponse[appdto.ListMembersResponse], errorbase.AppError)
		UpdateMemberRole(ctx context.Context, req *appdto.UpdateMemberRoleRequest) (*appdto.BaseResponse[appdto.MemberResponse], errorbase.AppError)
		RemoveMember(ctx context.Context, req *appdto.RemoveMemberRequest) (*appdto.BaseResponse[appdto.RemoveMemberResponse], errorbase.AppError)
		CreateInvite(ctx context.Context, req *appdto.CreateInviteRequest) (*appdto.BaseResponse[appdto.InviteResponse], errorbase.AppError)
		AcceptInvite(ctx context.Context, req *appdto.AcceptInviteRequest) (*appdto.BaseResponse[appdto.AcceptInviteResponse], errorbase.AppError)
		GeneratePresignedURLs(ctx context.Context, req *appdto.GeneratePresignedURLsRequest) (*appdto.BaseResponse[appdto.GeneratePresignedURLsResponse], errorbase.AppError)
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
		GenerateSprint(ctx context.Context, req *appdto.GenerateSprintRequest) (*appdto.BaseResponse[appdto.GenerateSprintResponse], errorbase.AppError)
		ConsumeAISprintGenerationResult(ctx context.Context) func(d rabbitmq.Delivery) rabbitmq.Action
		DeleteDraftSprint(ctx context.Context, req *appdto.DeleteSprintRequest) (*appdto.BaseResponse[appdto.DeleteSprintResponse], errorbase.AppError)
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
		GetUserInfo(ctx context.Context, req *appdto.UserInfoRequest) (*appdto.BaseResponse[appdto.UserInfoResponse], errorbase.AppError)
		NotificationConfiguration(ctx context.Context, req *appdto.ConfigureNotificationRequest) (*appdto.BaseResponse[appdto.ConfigureNotificationResponse], errorbase.AppError)
	}
)

func NewGroupUseCase(
	store istore.Store,
	validator *appvalidation.GroupValidator,
	authHelper *apphelper.AuthHelper,
	notificationHelper *apphelper.NotificationHelper,
	r2Client *storage.R2Client,
) GroupUseCase {
	return &groupUseCase{
		store:              store,
		groupRepo:          store.GroupRepository(),
		userRepo:           store.UserRepository(),
		validator:          validator,
		authHelper:         authHelper,
		notificationHelper: notificationHelper,
		r2Client:           r2Client,
	}
}

func NewSprintUseCase(
	store istore.Store,
	validator *appvalidation.SprintValidator,
	authHelper *apphelper.AuthHelper,
	cacheRepo icacherepository.CacheRepository,
	notificationHelper *apphelper.NotificationHelper,
	aiHelper *apphelper.AIHelper,
	logger log.LoggerV2,
) SprintUseCase {
	return &sprintUseCase{
		store:              store,
		sprintRepo:         store.SprintRepository(),
		workRepo:           store.WorkRepository(),
		userRepo:           store.UserRepository(),
		validator:          validator,
		authHelper:         authHelper,
		cacheRepo:          cacheRepo,
		sprintExportHelper: apphelper.NewSprintExportHelper(),
		groupRepo:          store.GroupRepository(),
		notificationHelper: notificationHelper,
		aiHelper:           aiHelper,
		logger:             logger,
	}
}

func NewWorkUseCase(
	store istore.Store,
	validator *appvalidation.WorkValidator,
	authHelper *apphelper.AuthHelper,
	notificationHelper *apphelper.NotificationHelper,
) WorkUseCase {
	return &workUseCase{
		store:              store,
		workRepo:           store.WorkRepository(),
		groupRepo:          store.GroupRepository(),
		validator:          validator,
		authHelper:         authHelper,
		notificationHelper: notificationHelper,
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
