package errordictionary

import (
	"fmt"
	coreerror "team_service/internal/domain/common/apperror"
)

const WORK_ERROR_DOMAIN = "err:work"

type WorkErrorInfo struct {
	ErrorInfo coreerror.ErrorInfo
	Id        string
}

var (
	ErrWorkNotFound = WorkErrorInfo{
		ErrorInfo: coreerror.ErrNotFound,
		Id:        fmt.Sprintf("%s:01", WORK_ERROR_DOMAIN),
	}
)
