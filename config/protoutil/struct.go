// Copyright (c) 2024 OrigAdmin. All rights reserved.

// Package protoutil provides helper functions for working with flexible configuration
// structures in Protobuf, such as google.protobuf.Struct and google.protobuf.Any.
// It simplifies the process of converting these generic types into strongly-typed
// Go structs derived from Protobuf messages.
package protoutil

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// UnmarshalTo unmarshals a generic *structpb.Struct into a specific
// strongly-typed protobuf message.
//
// It leverages the intermediate JSON representation of google.protobuf.Struct to
// provide the flexibility of a "normal" config file structure while allowing for
// strong typing within the application code.
//
// Usage:
//
//	var authzCfg authzv1.AuthzConfig
//	if err := customize.UnmarshalTo(middleware.Customize, &authzCfg); err != nil {
//	    // Handle error
//	}
//	// authzCfg is now populated with data from the config file.
func UnmarshalTo(s *structpb.Struct, dest proto.Message) error {
	// If the struct is nil, it means no configuration was provided.
	// This is a valid scenario, not an error.
	if s == nil {
		return nil
	}

	// Marshal the structpb.Struct into its canonical JSON representation.
	jsonBytes, err := protojson.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal struct to JSON: %w", err)
	}

	// Unmarshal the JSON bytes into the destination protobuf message.
	if err := protojson.Unmarshal(jsonBytes, dest); err != nil {
		return fmt.Errorf("failed to unmarshal JSON into %T: %w", dest, err)
	}

	return nil
}

// NewFromStruct creates a new instance of a strongly-typed protobuf message
// and unmarshals the generic *structpb.Struct into it.
//
// The generic type parameter 'T' must be the struct type of the protobuf message,
// and its pointer type (*T) must implement proto.Message.
//
// This function is useful when you want to directly receive the configured object
// as a return value, rather than passing a destination pointer.
//
// Usage:
//
//	authzCfg, err := customize.NewFromStruct[authzv1.AuthzConfig](middleware.Customize)
//	if err != nil {
//	    // Handle error
//	}
//	// authzCfg is now a new instance populated with data from the config file.
func NewFromStruct[T any, PT ProtoMessagePtr[T]](s *structpb.Struct) (PT, error) {
	// Create a new instance of the target type.
	// `new(T)` creates a pointer of type *T, which matches the return type PT.
	newVal := PT(new(T))

	// If no config is provided, return a new zero-valued instance of the target type.
	if s == nil {
		return newVal, nil
	}

	jsonBytes, err := protojson.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct to JSON: %w", err)
	}

	if err := protojson.Unmarshal(jsonBytes, newVal); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON into %T: %w", newVal, err)
	}

	return newVal, nil
}

// UnmarshalFromMap retrieves a specific configuration by key from a map of structs
// and unmarshals it into the destination protobuf message.
//
// Usage:
//
//	var authzCfg authzv1.AuthzConfig
//	err := customize.UnmarshalFromMap(myMap, "authz", &authzCfg)
func UnmarshalFromMap[T proto.Message](m map[string]*structpb.Struct, key string, dest T) error {
	if m == nil {
		return nil
	}
	s, ok := m[key]
	if !ok {
		return nil // Key not found is not an error, dest will remain zero-valued.
	}
	return UnmarshalTo(s, dest)
}

// NewFromStructMap retrieves a specific configuration by key from a map of structs,
// creates a new instance of the target type, and unmarshals the data into it.
//
// Usage:
//
//	authzCfg, err := customize.NewFromStructMap[authzv1.AuthzConfig](myMap, "authz")
func NewFromStructMap[T any, PT ProtoMessagePtr[T]](m map[string]*structpb.Struct, key string) (PT, error) {
	if m == nil {
		// If map is nil, return a new zero-valued instance.
		return PT(new(T)), nil
	}
	s, ok := m[key]
	if !ok {
		// If key is not found, also return a new zero-valued instance.
		return PT(new(T)), nil
	}
	// Reuse the existing NewTypedConfig logic.
	return NewFromStruct[T, PT](s)
}
