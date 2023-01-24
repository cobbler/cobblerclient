package raw

// Item general fields
type Item struct {
	// Internal XML-RPC handle
	Handle string

	// Item fields
	Parent            string            `mapstructure:"parent"`
	Depth             int               `mapstructure:"depth"          cobbler:"noupdate"`
	Children          []string          `mapstructure:"children"`
	CTime             float64           `mapstructure:"ctime"          cobbler:"noupdate"`
	MTime             float64           `mapstructure:"mtime"          cobbler:"noupdate"`
	Uid               string            `mapstructure:"uid"            cobbler:"noupdate"`
	Name              string            `mapstructure:"name"`
	Comment           string            `mapstructure:"comment"`
	KernelOptions     interface{}       `mapstructure:"kernel_options"`
	KernelOptionsPost interface{}       `mapstructure:"kernel_options_post"`
	AutoinstallMeta   interface{}       `mapstructure:"autoinstall_meta"`
	FetchableFiles    interface{}       `mapstructure:"fetchable_files"`
	BootFiles         map[string]string `mapstructure:"boot_files"`
	TemplateFiles     map[string]string `mapstructure:"template_files"`
	Owners            interface{}       `mapstructure:"owners"`
	MgmtClasses       interface{}       `mapstructure:"mgmt_classes"`
	MgmtParameters    interface{}       `mapstructure:"mgmt_parameters"`
}
