package item

type MgmtClass struct {
	Item `mapstructure:",squash"`

	// Mgmtclass specific fields
	IsDefiniton bool
	Params      map[string]interface{}
	ClassName   string
	Files       []string
	Packages    []string
}
