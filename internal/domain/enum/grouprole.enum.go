package enum

type GroupRole string

const (
	GroupRoleOwner   GroupRole = "owner"
	GroupRoleManager GroupRole = "manager"
	GroupRoleMember  GroupRole = "member"
	GroupRoleViewer  GroupRole = "viewer"
)

func (r GroupRole) IsValid() bool {
	switch r {
	case GroupRoleOwner,
		GroupRoleManager,
		GroupRoleMember,
		GroupRoleViewer:
		return true
	}
	return false
}

func (r GroupRole) Priority() int {
	switch r {
	case GroupRoleOwner:
		return 4
	case GroupRoleManager:
		return 3
	case GroupRoleMember:
		return 2
	case GroupRoleViewer:
		return 1
	default:
		return 0
	}
}

func (r GroupRole) HasPermission(required GroupRole) bool {
	return r.Priority() >= required.Priority()
}
