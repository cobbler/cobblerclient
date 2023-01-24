package item

type Resource struct {
	// Resource specific fields
	Action   ResourceAction `mapstructure:"action"`
	Mode     string         `mapstructure:"mode"`
	Owner    string         `mapstructure:"owner"`
	Group    string         `mapstructure:"group"`
	Path     string         `mapstructure:"path"`
	Template string         `mapstructure:"template"`
}
