package protocol

// ConfigEntry represents a named configuration item in an ordered sequence.
type ConfigEntry struct {
	Name  string
	Value any
}

// ConfigBlock is the standard interface for an ordered configuration block.
type ConfigBlock interface {
	GetActive() string
	GetDefault() any
	GetConfigs() []ConfigEntry
}

// ModuleConfig is the standardized output of an Extractor.
type ModuleConfig struct {
	Entries []ConfigEntry
	Active  string
}

// Extractor is a pure data retrieval function. 
// It does NOT need to perceive Scope.
type Extractor func(root any) (*ModuleConfig, error)

// Identifiable represents a configuration item that knows its own name.
type Identifiable interface {
	GetName() string
}

// Wrapper adapts arbitrary data into a ConfigBlock.
type Wrapper struct {
	Active  string
	Default any
	Configs []ConfigEntry
}

func (w *Wrapper) GetActive() string         { return w.Active }
func (w *Wrapper) GetDefault() any           { return w.Default }
func (w *Wrapper) GetConfigs() []ConfigEntry { return w.Configs }
