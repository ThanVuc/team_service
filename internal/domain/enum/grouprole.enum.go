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
