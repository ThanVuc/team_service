package utils

import (
	"math"
	appdto "team_service/internal/application/common/dto"
	"team_service/proto/common"
)

func ToPagination(pageQuery *common.PageQuery) appdto.Pagination {
	if pageQuery == nil {
		return appdto.Pagination{
			Page:  1,
			Limit: 10,
		}
	}

	page := int(pageQuery.Page)
	limit := int(pageQuery.PageSize)

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	if pageQuery.PageIgnore != nil && *pageQuery.PageIgnore {
		page = 1
		limit = 0
	}

	return appdto.Pagination{
		Page:  page,
		Limit: limit,
	}
}

func ToPageInfo(page, pageSize, totalItems int32) *common.PageInfo {
	totalPages := int32(math.Ceil(float64(totalItems) / float64(pageSize)))

	return &common.PageInfo{
		TotalItems: totalItems,
		Page:       page,
		TotalPages: totalPages,
		PageSize:   pageSize,
		HasPrev:    page > 1,
		HasNext:    page < totalPages,
	}
}
