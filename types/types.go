package types

// Encoder defines a function type for encoding a value into a byte slice.
type Encoder func(v any) ([]byte, error)

// GlobalDefaultKey is a constant string representing the global default key.
const GlobalDefaultKey = "default"
