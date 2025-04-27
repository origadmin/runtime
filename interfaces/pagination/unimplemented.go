/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package pagination implements the functions, types, and interfaces for the module.
package pagination

// UnimplementedRequester is a struct that implements the Requester interface with empty methods.
type UnimplementedRequester struct{}

// GetCurrent returns the current page number.
func (u UnimplementedRequester) GetCurrent() int32 {
	return 0
}

// GetPageSize returns the page size.
func (u UnimplementedRequester) GetPageSize() int32 {
	return 0
}

// GetPageToken returns the page token.
func (u UnimplementedRequester) GetPageToken() string {
	return ""
}

// GetOnlyCount returns a boolean indicating whether to only return the count.
func (u UnimplementedRequester) GetOnlyCount() bool {
	return false
}

// GetNoPaging returns a boolean indicating whether to disable paging.
func (u UnimplementedRequester) GetNoPaging() bool {
	return false
}

// GetOrderBy returns the order by fields.
func (u UnimplementedRequester) GetOrderBy() []string {
	return nil
}

// GetFieldMask returns the field mask.
func (u UnimplementedRequester) GetFieldMask() []string {
	return nil
}

// UnimplementedResponder is a struct that implements the Responder interface with empty methods.
type UnimplementedResponder struct{}

// GetSuccess returns a boolean indicating whether the request was successful.
func (u UnimplementedResponder) GetSuccess() bool {
	return false
}

// GetTotal returns the total number of items.
func (u UnimplementedResponder) GetTotal() int32 {
	return 0
}

// GetData returns the data.
func (u UnimplementedResponder) GetData() any {
	return nil
}

// GetCurrent returns the current page number.
func (u UnimplementedResponder) GetCurrent() int32 {
	return 0
}

// GetPageSize returns the page size.
func (u UnimplementedResponder) GetPageSize() int32 {
	return 0
}

// GetNextPageToken returns the next page token.
func (u UnimplementedResponder) GetNextPageToken() string {
	return ""
}

// GetExtra returns any extra data.
func (u UnimplementedResponder) GetExtra() any {
	return nil
}

// These variables ensure that UnimplementedRequester and UnimplementedResponder implement the Requester and Responder interfaces.
var _ Requester = (*UnimplementedRequester)(nil)
var _ Responder = (*UnimplementedResponder)(nil)
