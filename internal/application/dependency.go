package application

import (
	"team_service/internal/infrastructure"
)

type Dependency struct {
}

func NewDependency(infra *infrastructure.Dependency) *Dependency {
	return &Dependency{}
}
