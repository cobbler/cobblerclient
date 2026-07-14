package cobblerclient

// BuildisoOptions is a struct which describes the options one can set for the buildiso action of Cobbler.
type BuildisoOptions struct {
	Iso            string   `mapstructure:"iso" xmlrpc:"iso"`
	Profiles       []string `mapstructure:"profiles" xmlrpc:"profiles"`
	Systems        []string `mapstructure:"systems" xmlrpc:"systems"`
	BuildisoDir    string   `mapstructure:"buildisodir" xmlrpc:"buildisodir"`
	Distro         string   `mapstructure:"distro" xmlrpc:"distro"`
	Standalone     bool     `mapstructure:"standalone" xmlrpc:"standalone"`
	Airgapped      bool     `mapstructure:"airgapped" xmlrpc:"airgapped"`
	Source         string   `mapstructure:"source" xmlrpc:"source"`
	ExcludeDns     bool     `mapstructure:"exclude_dns" xmlrpc:"exclude_dns"`
	ExcludeSystems bool     `mapstructure:"exclude_systems" xmlrpc:"exclude_systems"`
	XorrisofsOpts  string   `mapstructure:"xorrisofs_opts" xmlrpc:"xorrisofs_opts"`
	Esp            string   `mapstructure:"esp" xmlrpc:"esp"`
}

// AclSetupOptions is a struct which describes the options one can set for the actlsetup action of Cobbler.
type AclSetupOptions struct {
	AddUser     string `mapstructure:"adduser" xmlrpc:"adduser"`
	AddGroup    string `mapstructure:"addgroup" xmlrpc:"addgroup"`
	RemoveUser  string `mapstructure:"removeuser" xmlrpc:"removeuser"`
	RemoveGroup string `mapstructure:"removegroup" xmlrpc:"removegroup"`
}

// ReplicateOptions is a struct which descibres the options one can set for the replicate action of Cobbler.
type ReplicateOptions struct {
	Master          string `mapstructure:"master" xmlrpc:"master"`
	Port            string `mapstructure:"port" xmlrpc:"port"`
	DistroPatterns  string `mapstructure:"distro_patterns" xmlrpc:"distro_patterns"`
	ProfilePatterns string `mapstructure:"profile_patterns" xmlrpc:"profile_patterns"`
	SystemPatterns  string `mapstructure:"system_patterns" xmlrpc:"system_patterns"`
	RepoPatterns    string `mapstructure:"repo_patterns" xmlrpc:"repo_patterns"`
	Imagepatterns   string `mapstructure:"image_patterns" xmlrpc:"image_patterns"`
	Prune           bool   `mapstructure:"prune" xmlrpc:"prune"`
	OmitData        bool   `mapstructure:"omit_data" xmlrpc:"omit_data"`
	SyncAll         bool   `mapstructure:"sync_all" xmlrpc:"sync_all"`
	UseSsl          bool   `mapstructure:"use_ssl" xmlrpc:"use_ssl"`
}

type BackgroundSyncOptions struct {
	Dhcp    bool `mapstructure:"dhcp" xmlrpc:"dhcp"`
	Dns     bool `mapstructure:"dns" xmlrpc:"dns"`
	Verbose bool `mapstructure:"verbose" xmlrpc:"verbose"`
}

type BackgroundSyncSystemsOptions struct {
	Systems []string `mapstructure:"systems" xmlrpc:"systems"`
	Verbose bool     `mapstructure:"verbose" xmlrpc:"verbose"`
}

type BackgroundImportOptions struct {
	Path            string `mapstructure:"path" xmlrpc:"path"`
	Name            string `mapstructure:"name" xmlrpc:"name"`
	AvailableAs     string `mapstructure:"available_as" xmlrpc:"available_as"`
	AutoinstallFile string `mapstructure:"autoinstall_file" xmlrpc:"autoinstall_file"`
	RsyncFlags      string `mapstructure:"rsync_flags" xmlrpc:"rsync_flags"`
	Arch            string `mapstructure:"arch" xmlrpc:"arch"`
	Breed           string `mapstructure:"breed" xmlrpc:"breed"`
	OsVersion       string `mapstructure:"os_version" xmlrpc:"os_version"`
}

type BackgroundReposyncOptions struct {
	Repos  []string `mapstructure:"repos" xmlrpc:"repos"`
	Only   string   `mapstructure:"only" xmlrpc:"only"`
	Nofail bool     `mapstructure:"nofail" xmlrpc:"nofail"`
	Tries  int      `mapstructure:"tries" xmlrpc:"tries"`
}

type BackgroundPowerSystemOptions struct {
	Systems []string `mapstructure:"systems" xmlrpc:"systems"`
	Power   string   `mapstructure:"power" xmlrpc:"power"`
}
