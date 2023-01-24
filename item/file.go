package item

// File is ...
type File struct {
	Item `mapstructure:",squash"`

	Resource `mapstructure:",squash"`

	// File specific fields
	IsDir bool `mapstructure:"is_dir"`
}
