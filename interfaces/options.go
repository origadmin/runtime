package interfaces

// Option used to store and retrieve option values
type Option interface {
	Value(key any) any
	With(key any, value any) Option
}
