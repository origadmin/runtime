package metadata

type Meta interface {
	Append(key string, values ...string)
	Values(key string) []string
	Get(key string) string
	Set(key string, value string)
	Clone() Meta
	GetAll() map[string][]string
}
