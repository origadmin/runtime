// Copyright (c) 2024 OrigAdmin. All rights reserved.

// Package customize provides helper functions for working with flexible configuration
// structures in Protobuf, such as google.protobuf.Struct and google.protobuf.Any.
// It simplifies the process of converting these generic types into strongly-typed
// Go structs derived from Protobuf messages.
package protobufutil

import "google.golang.org/protobuf/proto"

// ProtoMessagePtr is a generic constraint for a type that is a pointer to a struct T
// and also implements the proto.Message interface. This allows for creating generic
// functions that can work with any pointer to a protobuf message struct.
type ProtoMessagePtr[T any] interface {
	*T
	proto.Message
}
