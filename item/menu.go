package item

// Menu is ...
type Menu struct {
	Item `mapstructure:",squash"`

	// Menu specific fields
	DisplayName string `mapstructure:"display_name"`
}
