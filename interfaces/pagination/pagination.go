/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package pagination implements the functions, types, and interfaces for the module.
package pagination

// CurrentGetter is the interface for getting the current page number.
type CurrentGetter interface {
	GetCurrent() int32
}

// PageSizeGetter is the interface for getting the page size.
type PageSizeGetter interface {
	GetPageSize() int32
}

// PageTokenGetter is the interface for getting the page token.
type PageTokenGetter interface {
	GetPageToken() string
}

// OnlyCountGetter is the interface for checking if only count is needed.
type OnlyCountGetter interface {
	GetOnlyCount() bool
}

// NoPagingGetter is the interface for checking if no paging is needed.
type NoPagingGetter interface {
	GetNoPaging() bool
}

// OrderByGetter is the interface for getting the order by fields.
type OrderByGetter interface {
	GetOrderBy() []string
}

// FieldMaskGetter is the interface for getting the field mask.
type FieldMaskGetter interface {
	GetFieldMask() []string
}

// SuccessGetter is the interface for getting the success status.
type SuccessGetter interface {
	GetSuccess() bool
}

// TotalGetter is the interface for getting the total count.
type TotalGetter interface {
	GetTotal() int32
}

// DataGetter is the interface for getting the data.
type DataGetter interface {
	GetData() any
}

// NextPageTokenGetter is the interface for getting the next page token.
type NextPageTokenGetter interface {
	GetNextPageToken() string
}

// ExtraGetter is the interface for getting the extra data.
type ExtraGetter interface {
	GetExtra() any
}

// Requester is the request interface for the module.
type Requester interface {
	CurrentGetter
	PageSizeGetter
	PageTokenGetter
	OnlyCountGetter
	NoPagingGetter
	OrderByGetter
	FieldMaskGetter
}

// Responder is the response interface for the module.
type Responder interface {
	SuccessGetter
	TotalGetter
	DataGetter
	CurrentGetter
	PageSizeGetter
	NextPageTokenGetter
	ExtraGetter
}

type PageRequest interface {
	CurrentGetter
	PageSizeGetter
	PageTokenGetter
}

type FilterRequest interface {
	OrderByGetter
	FieldMaskGetter
}

type ControlRequest interface {
	OnlyCountGetter
	NoPagingGetter
}

type PageResponse interface {
	CurrentGetter
	PageSizeGetter
	NextPageTokenGetter
}

type DataResponse interface {
	TotalGetter
	DataGetter
	ExtraGetter
}
