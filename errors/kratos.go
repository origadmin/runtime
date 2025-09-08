/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package errors provides a centralized hub for handling, converting, and rendering errors.
// It provides custom error handlers/encoders for HTTP and gRPC that integrate with the
// Kratos ecosystem while providing centralized logging and error conversion.
package errors

//go:generate adptool .
//go:adapter:package github.com/go-kratos/kratos/v2/errors kerrors
