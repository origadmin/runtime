package options

// Context is used to store and retrieve option values, similar to context.Context.
type Context interface {
	Value(key any) any
	With(key any, value any) Context
}

// Option is a function that applies a configuration to a Context.
type Option func(c Context) Context
