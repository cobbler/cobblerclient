package item

// Profile is a created profile.
type Profile struct {
	Item `mapstructure:",squash"`

	// These are internal fields and cannot be modified.
	ReposEnabled bool `mapstructure:"repos_enabled"          cobbler:"noupdate"`

	Autoinstall       string   `mapstructure:"autoinstall"`
	DHCPTag           string   `mapstructure:"dhcp_tag"`
	Distro            string   `mapstructure:"distro"`
	EnableGPXE        bool     `mapstructure:"enable_gpxe"`
	EnableMenu        bool     `mapstructure:"enable_menu"`
	NameServers       []string `mapstructure:"name_servers"`
	NameServersSearch []string `mapstructure:"name_servers_search"`
	NextServerv4      string   `mapstructure:"next_server_v4"`
	NextServerv6      string   `mapstructure:"next_server_v6"`
	Parent            string   `mapstructure:"parent"`
	Proxy             string   `mapstructure:"proxy"`
	Repos             []string `mapstructure:"repos"`
	Server            string   `mapstructure:"server"`
	VirtAutoBoot      string   `mapstructure:"virt_auto_boot"`
	VirtBridge        string   `mapstructure:"virt_bridge"`
	VirtCPUs          string   `mapstructure:"virt_cpus"`
	VirtDiskDriver    string   `mapstructure:"virt_disk_driver"`
	VirtFileSize      string   `mapstructure:"virt_file_size"`
	VirtPath          string   `mapstructure:"virt_path"`
	VirtRAM           string   `mapstructure:"virt_ram"`
	VirtType          string   `mapstructure:"virt_type"`
}
