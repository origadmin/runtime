/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

// UnimplementedClaims is a struct that implements the Claims interface
type UnimplementedClaims struct {
}

func (u UnimplementedClaims) GetJWTID() string {
	return ""
}

// GetSubject returns an empty string
func (u UnimplementedClaims) GetSubject() string {
	return ""
}

// GetIssuer returns an empty string
func (u UnimplementedClaims) GetIssuer() string {
	return ""
}

// GetAudience returns an empty slice
func (u UnimplementedClaims) GetAudience() []string {
	return []string{}
}

// GetExpiration returns the current time
func (u UnimplementedClaims) GetExpiration() int64 {
	return 0
}

// GetNotBefore returns the current time
func (u UnimplementedClaims) GetNotBefore() int64 {
	return 0
}

// GetIssuedAt returns the current time
func (u UnimplementedClaims) GetIssuedAt() int64 {
	return 0
}

// GetID returns an empty string
func (u UnimplementedClaims) GetID() string {
	return ""
}

// GetScopes returns an empty map
func (u UnimplementedClaims) GetScopes() map[string]bool {
	return make(map[string]bool)
}

type UnimplementedPolicy struct {
}

func (u UnimplementedPolicy) GetRoles() []string {
	return nil
}

func (u UnimplementedPolicy) GetPermissions() []string {
	return nil
}

func (u UnimplementedPolicy) GetSubject() string {
	return ""
}

func (u UnimplementedPolicy) GetObject() string {
	return ""
}

func (u UnimplementedPolicy) GetAction() string {
	return ""
}

func (u UnimplementedPolicy) GetDomain() string {
	return ""
}

var _ Claims = (*UnimplementedClaims)(nil)
var _ Policy = (*UnimplementedPolicy)(nil)
