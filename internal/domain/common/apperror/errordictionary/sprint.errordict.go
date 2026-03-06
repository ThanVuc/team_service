package errordictionary

import (
	"fmt"
	coreerror "team_service/internal/domain/common/apperror"
)

const SPRINT_ERROR_DOMAIN = "err:sprint"

type SprintErrorInfo struct {
	ErrorInfo coreerror.ErrorInfo
	Id        string
}

var (
	ErrSprintNotFound = SprintErrorInfo{
		ErrorInfo: coreerror.ErrNotFound,
		Id:        fmt.Sprintf("%s:01", SPRINT_ERROR_DOMAIN),
	}
)
