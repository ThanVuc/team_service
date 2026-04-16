package apphelper

import (
	"strings"
	appdto "team_service/internal/application/common/dto"
	"team_service/internal/domain/enum"
)

func CollectMemberIDsByRoles(members *appdto.ListMembersResponse, roles ...enum.GroupRole) []string {
	if members == nil || len(members.Members) == 0 {
		return []string{}
	}

	roleFilter := make(map[enum.GroupRole]struct{}, len(roles))
	for _, role := range roles {
		roleFilter[role] = struct{}{}
	}

	receivers := make([]string, 0, len(members.Members))
	for _, member := range members.Members {
		if _, ok := roleFilter[member.Role]; !ok {
			continue
		}

		receivers = append(receivers, member.ID)
	}

	return UniqueIDs(receivers...)
}

func CollectAllMemberIDs(members *appdto.ListMembersResponse) []string {
	if members == nil || len(members.Members) == 0 {
		return []string{}
	}

	receivers := make([]string, 0, len(members.Members))
	for _, member := range members.Members {
		receivers = append(receivers, member.ID)
	}

	return UniqueIDs(receivers...)
}

func CollectDiscussionParticipantIDs(work *appdto.WorkResponse) []string {
	if work == nil || work.Comments == nil || len(work.Comments.Comments) == 0 {
		return []string{}
	}

	participantIDs := make([]string, 0, len(work.Comments.Comments))
	for _, comment := range work.Comments.Comments {
		participantIDs = append(participantIDs, comment.Creator.ID)
	}

	return UniqueIDs(participantIDs...)
}

func UniqueIDs(ids ...string) []string {
	seen := make(map[string]struct{}, len(ids))
	result := make([]string, 0, len(ids))

	for _, id := range ids {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}

		if _, ok := seen[trimmed]; ok {
			continue
		}

		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}

	return result
}

func ExcludeID(ids []string, excludedID string) []string {
	excludedID = strings.TrimSpace(excludedID)
	if excludedID == "" || len(ids) == 0 {
		return ids
	}

	filtered := make([]string, 0, len(ids))
	for _, id := range ids {
		if strings.TrimSpace(id) == excludedID {
			continue
		}
		filtered = append(filtered, id)
	}

	return filtered
}