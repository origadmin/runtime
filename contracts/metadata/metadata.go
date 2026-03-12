package metadata

type Metadata interface {
	Append(key string, values ...string)
	Values(key string) []string
	Get(key string) string
	Set(key string, value string)
	Clone() Metadata
	GetAll() map[string][]string
}
