package apphelper

import (
	"context"
	"fmt"
	"strings"
	"team_service/internal/infrastructure/share/utils"
)

func ResolveNotificationOrigin(ctx context.Context) string {
	origin := strings.TrimSpace(utils.GetOriginFromIncomingContext(ctx))
	return strings.TrimRight(origin, "/")
}

func BuildMembersTabLink(ctx context.Context, groupID string) string {
	return fmt.Sprintf("%s/te/group/%s?tab=members", ResolveNotificationOrigin(ctx), strings.TrimSpace(groupID))
}

func BuildSprintsTabLink(ctx context.Context, groupID string) string {
	return fmt.Sprintf("%s/te/group/%s?tab=sprints", ResolveNotificationOrigin(ctx), strings.TrimSpace(groupID))
}

func BuildSprintWorkboardLink(ctx context.Context, groupID string, sprintID string) string {
	return fmt.Sprintf(
		"%s/te/group/%s?tab=workboard&sprint_id=%s",
		ResolveNotificationOrigin(ctx),
		strings.TrimSpace(groupID),
		strings.TrimSpace(sprintID),
	)
}

func BuildWorkUpdateLink(ctx context.Context, groupID string, sprintID string, workID string) string {
	return fmt.Sprintf(
		"%s/te/group/%s?tab=workboard&sprint_id=%s&mode=Update&id=%s",
		ResolveNotificationOrigin(ctx),
		strings.TrimSpace(groupID),
		strings.TrimSpace(sprintID),
		strings.TrimSpace(workID),
	)
}
