package query

import (
	"github.com/origadmin/runtime/interfaces/query"
)

// NewPageRequest creates a new PageRequest with default values.
func NewPageRequest(page, pageSize int, sortBy ...query.SortOrder) query.PageRequest {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // Default page size
	}
	return query.PageRequest{
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
	}
}

// NewPageResponse creates a new PageResponse.
func NewPageResponse[T any](content []T, totalElements int64, pageRequest query.PageRequest) query.PageResponse[T] {
	totalPages := 0
	if pageRequest.PageSize > 0 {
		totalPages = int((totalElements + int64(pageRequest.PageSize) - 1) / int64(pageRequest.PageSize))
	}
	return query.PageResponse[T]{
		Content:       content,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Page:          pageRequest.Page,
		PageSize:      pageRequest.PageSize,
	}
}
