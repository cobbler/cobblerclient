package raw

// Distro is a created distro.
// Get the fields from cobbler/items/distro.py
type Distro struct {
	Item `mapstructure:",squash"`

	// These are internal fields and cannot be modified.
	SourceRepos    []string    `mapstructure:"source_repos"    cobbler:"noupdate"`
	TreeBuildTime  string      `mapstructure:"tree_build_time" cobbler:"noupdate"`
	Arch           string      `mapstructure:"arch"`
	BootFiles      []string    `mapstructure:"boot_files"`
	BootLoaders    interface{} `mapstructure:"boot_loaders"`
	Breed          string      `mapstructure:"breed"`
	FetchableFiles []string    `mapstructure:"fetchable_files"`
	Initrd         string      `mapstructure:"initrd"`
	Kernel         string      `mapstructure:"kernel"`
	MGMTClasses    []string    `mapstructure:"mgmt_classes"`
	OSVersion      string      `mapstructure:"os_version"`
	TemplateFiles  []string    `mapstructure:"template_files"`
}
