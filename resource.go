package cobblerclient

type ResourceAction int64

const (
	raCREATE = iota
	raREMOVE
)

func (r ResourceAction) String() string {
	switch r {
	case raCREATE:
		return "create"
	case raREMOVE:
		return "remove"
	}
	return "unknown"
}

// Resource is an abstract item type that cannot be directly instantiated.
// Get the fields from cobbler/items/resource.py
type Resource struct {
	Item `mapstructure:",squash"`

	// Resource specific attributes
	// Action   ResourceAction `mapstructure:"action"`
	Action   string `mapstructure:"action"`
	Mode     string `mapstructure:"mode"`
	Owner    string `mapstructure:"owner"`
	Group    string `mapstructure:"group"`
	Path     string `mapstructure:"path"`
	Template string `mapstructure:"template"`
}

func NewResource() Resource {
	return Resource{
		Item:   NewItem(),
		Action: "create",
	}
}
