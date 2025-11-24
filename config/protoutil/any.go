// Copyright (c) 2024 OrigAdmin. All rights reserved.

// Package protoutil provides utility functions for working with protobuf Any messages.
// It includes methods for packing and unpacking protobuf messages into and from Any types,
// as well as creating new instances of strongly-typed messages from Any values.
package protoutil

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// Pack packs the provided protobuf message into an anypb.Any.
// It returns an error if the packing fails. This is a convenient
// wrapper around anypb.New().
//
// Usage:
//
//	anyValue, err := customize.Pack(&myCfg)
func Pack(msg proto.Message) (*anypb.Any, error) {
	anyMsg, err := anypb.New(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to pack message into Any: %w", err)
	}
	return anyMsg, nil
}

// UnpackTo unpacks the configuration from an anypb.Any into the destination
// message 'dest'. It returns an error if the types do not match or unpacking fails.
func UnpackTo(any *anypb.Any, dest proto.Message) error {
	if any == nil {
		return nil // It's not an error if there's nothing to unpack.
	}
	if err := any.UnmarshalTo(dest); err != nil {
		return fmt.Errorf("failed to unmarshal Any to %T: %w", dest, err)
	}
	return nil
}

// NewFromAny creates a new instance of a strongly-typed protobuf message
// and unmarshals the content of an anypb.Any into it.
// This is the counterpart to NewTypedConfig for the Any type.
func NewFromAny[T any, PT ProtoMessagePtr[T]](any *anypb.Any) (PT, error) {
	newVal := PT(new(T))
	if any == nil {
		return newVal, nil // Return a zero-valued instance if input is nil.
	}
	if err := any.UnmarshalTo(newVal); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Any to new %T: %w", newVal, err)
	}
	return newVal, nil
}

// UnpackFromMap retrieves a specific configuration by key from a map of Any messages
// and unmarshals it into the destination protobuf message.
//
// Usage:
//
//	var authzCfg authzv1.AuthzConfig
//	err := customize.UnpackFromMap(myMap, "authz", &authzCfg)
func UnpackFromMap(m map[string]*anypb.Any, key string, dest proto.Message) error {
	if m == nil {
		return nil
	}
	any, ok := m[key]
	if !ok {
		return nil // Key not found is not an error.
	}
	return UnpackTo(any, dest)
}

// NewFromAnyMap retrieves a specific configuration by key from a map of Any messages,
// creates a new instance of the target type, and unmarshals the data into it.
//
// Usage:
//
//	authzCfg, err := customize.NewFromAnyMap[authzv1.AuthzConfig](myMap, "authz")
func NewFromAnyMap[T any, PT ProtoMessagePtr[T]](m map[string]*anypb.Any, key string) (PT, error) {
	if m == nil {
		return PT(new(T)), nil
	}
	any, ok := m[key]
	if !ok {
		return PT(new(T)), nil
	}
	return NewFromAny[T, PT](any)
}
