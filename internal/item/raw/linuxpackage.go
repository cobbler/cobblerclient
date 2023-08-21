package raw

type Package struct {
	Item `mapstructure:",squash"`

	// Package specific attributes
	Installer string `mapstructure:"installer"`
	Version   string `mapstructure:"version"`
}
