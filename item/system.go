package item

// System is a created system.
type System struct {
	Item `mapstructure:",squash"`
	// These are internal fields and cannot be modified.
	IPv6Autoconfiguration bool                   `mapstructure:"ipv6_autoconfiguration" cobbler:"noupdate"`
	ReposEnabled          bool                   `mapstructure:"repos_enabled"          cobbler:"noupdate"`
	Autoinstall           string                 `mapstructure:"autoinstall"`
	BootFiles             string                 `mapstructure:"boot_files"`
	BootLoaders           []string               `mapstructure:"boot_loaders"`
	EnableGPXE            bool                   `mapstructure:"enable_gpxe"`
	FetchableFiles        []string               `mapstructure:"fetchable_files"`
	Gateway               string                 `mapstructure:"gateway"`
	Hostname              string                 `mapstructure:"hostname"`
	Image                 string                 `mapstructure:"image"`
	Interfaces            map[string]interface{} `mapstructure:"interfaces" cobbler:"noupdate"`
	IPv6DefaultDevice     string                 `mapstructure:"ipv6_default_device"`
	KernelOptions         []string               `mapstructure:"kernel_options"`
	KernelOptionsPost     []string               `mapstructure:"kernel_options_post"`
	MGMTClasses           []string               `mapstructure:"mgmt_classes"`
	MGMTParameters        string                 `mapstructure:"mgmt_parameters"`
	NameServers           []string               `mapstructure:"name_servers"`
	NameServersSearch     []string               `mapstructure:"name_servers_search"`
	NetbootEnabled        bool                   `mapstructure:"netboot_enabled"`
	NextServerv4          string                 `mapstructure:"next_server_v4"`
	NextServerv6          string                 `mapstructure:"next_server_v6"`
	PowerAddress          string                 `mapstructure:"power_address"`
	PowerID               string                 `mapstructure:"power_id"`
	PowerPass             string                 `mapstructure:"power_pass"`
	PowerType             string                 `mapstructure:"power_type"`
	PowerUser             string                 `mapstructure:"power_user"`
	Profile               string                 `mapstructure:"profile"`
	Proxy                 string                 `mapstructure:"proxy"`
	Status                string                 `mapstructure:"status"`
	TemplateFiles         []string               `mapstructure:"template_files"`
	VirtAutoBoot          string                 `mapstructure:"virt_auto_boot"`
	VirtCPUs              string                 `mapstructure:"virt_cpus"`
	VirtDiskDriver        string                 `mapstructure:"virt_disk_driver"`
	VirtFileSize          string                 `mapstructure:"virt_file_size"`
	VirtPath              string                 `mapstructure:"virt_path"`
	VirtPXEBoot           int                    `mapstructure:"virt_pxe_boot"`
	VirtRAM               string                 `mapstructure:"virt_ram"`
	VirtType              string                 `mapstructure:"virt_type"`
}
