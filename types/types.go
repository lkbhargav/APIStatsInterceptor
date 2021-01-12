package types

// Option => defines different option types
type Option string

const (
	// Comma => depicts a comma (,)
	Comma Option = "COMMA"
	// Percent => depicts a percent (%)
	Percent Option = "PERCENT"
	// Data => adds an appropriate suffix (KB, MB, etc.)
	Data Option = "DATA"
	// Prefix => adds a prefix to the value
	Prefix Option = "PREFIX"
	// Suffix => adds a suffix to the value
	Suffix Option = "SUFFIX"
	// None => just a nil
	None Option = "NONE"
)

// Set => gets the set information
type Set struct {
	Name        string   `json:"name"`
	Path        []string `json:"path"`
	Option      Option   `json:"option"`
	OptionalVal string   `json:"optionalValue"`
}
