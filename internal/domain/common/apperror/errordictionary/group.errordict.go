package errordictionary

import (
	"fmt"
	coreerror "team_service/internal/domain/common/apperror"
)

const GROUP_ERROR_DOMAIN = "err:group"

type GroupErrorInfo struct {
	ErrorInfo coreerror.ErrorInfo
	Id        string
}

func (e GroupErrorInfo) Error() coreerror.ErrorInfo {
	return e.ErrorInfo
}

var (
	ErrGroupNotFound = GroupErrorInfo{
		ErrorInfo: coreerror.ErrNotFound,
		Id:        fmt.Sprintf("%s:01", GROUP_ERROR_DOMAIN),
	}

	ErrGroupConflict = GroupErrorInfo{
		ErrorInfo: coreerror.ErrConflict,
		Id:        fmt.Sprintf("%s:02", GROUP_ERROR_DOMAIN),
	}

	ErrGroupBadRequest = GroupErrorInfo{
		ErrorInfo: coreerror.ErrBadRequest,
		Id:        fmt.Sprintf("%s:03", GROUP_ERROR_DOMAIN),
	}
)
