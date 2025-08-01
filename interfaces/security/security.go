/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.

// Package security is a package that provides security-related functions and types.
package security

// TokenSource type is defined to represent the origin of the token.
//
//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer@latest -type=TokenSource -trimprefix=TokenSource -output=security_string.go
type TokenSource int

// TokenSource constants represent the different types of context.
const (
	// TokenSourceContext represents the token source for the context.
	TokenSourceContext TokenSource = iota
	// TokenSourceHeader represents the token source for the header, if you don't know server or client
	TokenSourceHeader
	// TokenSourceClientHeader represents the token source for the header.
	TokenSourceClientHeader
	// TokenSourceServerHeader represents the token source for the header.
	TokenSourceServerHeader
	// TokenSourceMetadata represents the token source for the metadata, if you don't know server or client.
	TokenSourceMetadata
	// TokenSourceMetadataClient represents the token source for the metadata.
	TokenSourceMetadataClient
	// TokenSourceMetadataServer represents the token source for the metadata.
	TokenSourceMetadataServer
	// TokenSourceQueryParameter represents the token source for the query.
	TokenSourceQueryParameter
	// TokenSourceCookie represents the token source for the cookie.
	TokenSourceCookie
	// TokenSourceURLParameter represents the token source for the parameter.
	TokenSourceURLParameter
	// TokenSourceForm represents the token source for the form.
	TokenSourceForm
	// TokenSourceRequestBody represents the token source for the body.
	TokenSourceRequestBody
	// TokenSourceSession represents the token source for the session.
	TokenSourceSession
	// TokenSourceUnknown represents an unknown token source.
	TokenSourceUnknown
)

const (
	// HeaderAuthorize is the name of the authorization header.
	HeaderAuthorize = "Authorization"
	// HeaderContentType is the name of the content type header.
	HeaderContentType = "Content-Type"
	// HeaderContentLength is the name of the content length header.
	HeaderContentLength = "Content-Length"
	// HeaderUserAgent is the name of the user agent header.
	HeaderUserAgent = "User-Agent"
	// HeaderReferer is the name of the referer header.
	HeaderReferer = "Referer"
	// HeaderOrigin is the name of the origin header.
	HeaderOrigin = "Origin"
)

// Security represents the security interface.
type Security interface {
	Authenticator
	Authorizer
}
