package interfaces

// OptionValue used to store and retrieve option values
type OptionValue interface {
	Option(key any) any
	WithOption(key any, value any) OptionValue
}

// OptionSet implement the optionvalue interface
type OptionSet struct {
	parent OptionValue
	key    any
	value  any
}

func (c *OptionSet) Option(key any) any {
    if c == nil {
        return nil
    }
    if c.key != nil && c.key == key {
        return c.value
    }
    if c.parent != nil {
        return c.parent.Option(key)
    }
    return nil
}

func (c *OptionSet) WithOption(key any, value any) OptionValue {
    if key == nil {
        panic("key cannot be nil")
    }
    return &OptionSet{
        parent: c,
        key:    key,
        value:  value,
    }
}

func DefaultOptions() OptionValue {
	return &OptionSet{}
}
