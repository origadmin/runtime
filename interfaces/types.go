package interfaces

type Encoder func(v any) ([]byte, error)

const GlobalDefaultKey = "default"
