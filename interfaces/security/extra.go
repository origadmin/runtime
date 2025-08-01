/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

type Extra interface {
	// GetExtra returns the extra data as a map of strings
	GetExtra() map[string]string
	// Get returns the value associated with the given key
	Get(key string) (string, bool)
	// Set sets the value associated with the given key
	Set(key string, value string)
}

// ExtraData is an interface that defines methods for handling extra data associated with the security claims
type ExtraData interface {
	Extra
	// GetClaims returns the Claims object associated with the extra data,if Claims exists
	GetClaims() (Claims, bool)
	// HasClaims returns true if the extra data contains a Claims object
	HasClaims() bool
	// GetPolicy returns the Policy object associated with the extra data,if Policy exists
	GetPolicy() (Policy, bool)
	// HasPolicy returns true if the extra data contains a Policy object
	HasPolicy() bool
}

type extra map[string]string

func (e extra) GetExtra() map[string]string {
	return e
}

func (e extra) Get(key string) (string, bool) {
	ev, ok := e[key]
	return ev, ok
}

func (e extra) Set(key string, value string) {
	e[key] = value
}

type extraData struct {
	Claims Claims
	Policy Policy
	Extra  Extra
}

func (e extraData) HasClaims() bool {
	return e.Claims != nil
}

func (e extraData) HasPolicy() bool {
	return e.Policy != nil
}

func (e extraData) GetClaims() (Claims, bool) {
	return e.Claims, e.Claims != nil
}

func (e extraData) GetPolicy() (Policy, bool) {
	return e.Policy, e.Policy != nil
}

func (e extraData) HasExtra() bool {
	return e.Extra != nil
}
func (e extraData) GetExtra() map[string]string {
	if e.Extra == nil {
		return make(map[string]string)
	}
	return e.Extra.GetExtra()
}

func (e extraData) Get(key string) (string, bool) {
	if e.Extra == nil {
		return "", false
	}
	return e.Extra.Get(key)
}

func (e extraData) Set(key string, value string) {
	if e.Extra == nil {
		e.Extra = make(extra)
	}
	e.Extra.Set(key, value)
}

// ExtraObject retrieves the ExtraData object from a Policy if it implements the ExtraData interface
func ExtraObject(extra any) (Extra, bool) {
	if ex, ok := extra.(Extra); ok {
		return ex, true
	}
	return nil, false
}

func DataWithExtra(claims Claims, policy Policy, ext map[string]string) ExtraData {
	return &extraData{
		Claims: claims,
		Policy: policy,
		Extra:  (extra)(ext),
	}
}

func Claims2Extra(claims Claims, ext map[string]string) ExtraData {
	return &extraData{
		Claims: claims,
		Extra:  (extra)(ext),
	}
}

func Policy2Extra(policy Policy, ext map[string]string) ExtraData {
	return &extraData{
		Policy: policy,
		Extra:  (extra)(ext),
	}
}

func WithExtraData(ext map[string]string) ExtraData {
	return &extraData{
		Extra: (extra)(ext),
	}
}

func ClaimsWithExtra(claims Claims, ext map[string]string) Claims {
	if v, ok := claims.(ExtraClaims); ok {
		v.Extra = (extra)(ext)
		return v
	}
	return &ExtraClaims{
		Claims: claims,
		Extra:  (extra)(ext),
	}
}

func PolicyWithExtra(policy Policy, ext map[string]string) Policy {
	if v, ok := policy.(ExtraPolicy); ok {
		v.Extra = (extra)(ext)
		return v
	}
	return &ExtraPolicy{
		Policy: policy,
		Extra:  (extra)(ext),
	}
}

var _ ExtraData = (*extraData)(nil)
