package irepository

type RepositoryContainer struct {
	GroupRepository  *GroupRepository
	SprintRepository *SprintRepository
	WorkRepository   *WorkRepository
	InviteRepository *InviteRepository
}

func (c *RepositoryContainer) GetGroupRepository() *GroupRepository {
	return c.GroupRepository
}

func (c *RepositoryContainer) GetSprintRepository() *SprintRepository {
	return c.SprintRepository
}

func (c *RepositoryContainer) GetWorkRepository() *WorkRepository {
	return c.WorkRepository
}

func (c *RepositoryContainer) GetInviteRepository() *InviteRepository {
	return c.InviteRepository
}
